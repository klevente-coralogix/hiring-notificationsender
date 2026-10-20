//! The third-party API that messages are sent to, plus a mock for tests.

use std::collections::VecDeque;
use std::sync::Mutex;
use std::time::Duration;

/// Error returned by the third-party API. Mirrors its HTTP error responses.
#[derive(Debug, Clone, PartialEq, Eq, thiserror::Error)]
pub enum ApiError {
    /// 400: the message is malformed.
    #[error("400 Bad Request")]
    BadRequest,
    /// 401: our credentials are missing or invalid.
    #[error("401 Unauthorized")]
    Unauthorized,
    /// 413: the message is too large.
    #[error("413 Payload Too Large")]
    PayloadTooLarge,
    /// 429: we are sending too fast. `retry_after` is how long the API asks us
    /// to wait before sending again, if it says so.
    #[error("429 Too Many Requests (retry after: {retry_after:?})")]
    TooManyRequests { retry_after: Option<Duration> },
    /// 500: the API failed while processing the message.
    #[error("500 Internal Server Error")]
    InternalServerError,
    /// 503: the API is temporarily down or overloaded.
    #[error("503 Service Unavailable")]
    ServiceUnavailable,
    /// 504: an upstream of the API timed out. The message may or may not have been delivered.
    #[error("504 Gateway Timeout")]
    GatewayTimeout,
}

/// The third-party API. Sends a single message.
pub trait ThirdPartyApi: Send + Sync {
    fn post(&self, message: &str) -> impl Future<Output = Result<(), ApiError>> + Send;
}

/// A configurable fake of the third-party API.
///
/// Every call waits for the configured delay, then responds with the next
/// scripted response (see [`MockApi::set_responses`]), or with success once
/// those run out. Messages that were posted successfully are recorded in order.
pub struct MockApi {
    state: Mutex<State>,
}

struct State {
    delay: Duration,
    responses: VecDeque<Result<(), ApiError>>,
    posted: Vec<String>,
}

impl MockApi {
    /// Creates a mock that responds immediately with success.
    pub fn new() -> Self {
        Self {
            state: Mutex::new(State {
                delay: Duration::ZERO,
                responses: VecDeque::new(),
                posted: Vec::new(),
            }),
        }
    }

    /// Sets how long each `post` call waits before responding.
    pub fn set_delay(&self, delay: Duration) {
        self.state.lock().unwrap().delay = delay;
    }

    /// Scripts the responses of the next `post` calls, one per call, in order.
    /// Once they run out, calls respond with success.
    /// A scripted response is used up even if the call is cancelled (e.g. by a timeout).
    pub fn set_responses(&self, responses: impl IntoIterator<Item = Result<(), ApiError>>) {
        self.state.lock().unwrap().responses = responses.into_iter().collect();
    }

    /// Returns the messages that were posted successfully, in order.
    pub fn posted(&self) -> Vec<String> {
        self.state.lock().unwrap().posted.clone()
    }
}

impl Default for MockApi {
    fn default() -> Self {
        Self::new()
    }
}

impl ThirdPartyApi for MockApi {
    async fn post(&self, message: &str) -> Result<(), ApiError> {
        let (delay, response) = {
            let mut state = self.state.lock().unwrap();
            (state.delay, state.responses.pop_front().unwrap_or(Ok(())))
        };
        tokio::time::sleep(delay).await;

        if response.is_ok() {
            let mut state = self.state.lock().unwrap();
            state.posted.push(message.to_string());
        }
        response
    }
}
