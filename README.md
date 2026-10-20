# Notification Sender – Hiring Exercise

## Task

Implement a component for sending messages to a third-party API reliably.

Users can enqueue messages; the component processes them in order and deals with the third-party API gracefully. Users can also discard a message that is still waiting in the queue.

## Requirements

- Enqueueing a message returns immediately; it does not wait for the message to be sent.
- Messages are sent in the order they were enqueued.
- A message that has not been sent yet can be discarded.

## Provided

Each sub-project contains a `thirdparty` module with:

- The third-party API contract: a single `post` method that sends one message and either succeeds or fails. Its errors mirror the HTTP error responses of a typical web API (400, 401, 413, 429, 500, 503, 504).
- A mock implementation for tests. You can configure how long it waits before responding and whether it responds with success or an error (a scripted sequence of responses, one per call; success once the script runs out), and you can ask it which messages it has received.

## Deliverables

- A library (or binary with a library interface) in this repository, under the language sub-project you are working in.
- A brief explanation in this README of your design decisions.

## Sub-projects

### Rust

Located in [`rust/`](rust/). Uses the Rust standard library, [tokio](https://tokio.rs/) and [thiserror](https://docs.rs/thiserror).

```rust
use std::time::Duration;
use hiring_notificationsender::thirdparty::{ApiError, MockApi, ThirdPartyApi};

let api = MockApi::new();
api.set_delay(Duration::from_millis(100));
// The next two calls fail; every call after that succeeds.
api.set_responses([
    Err(ApiError::ServiceUnavailable),
    Err(ApiError::TooManyRequests { retry_after: Some(Duration::from_secs(1)) }),
]);

assert!(api.post("a").await.is_err());
assert!(api.post("b").await.is_err());
assert!(api.post("c").await.is_ok());
assert_eq!(api.posted(), ["c"]); // only successful posts are recorded
```

```sh
cd rust
cargo build
cargo test
cargo clippy
cargo fmt
cargo run
```

### Go

Located in [`go/`](go/). Uses the Go standard library only.

```go
api := thirdparty.NewMockAPI()
api.SetDelay(100 * time.Millisecond)
// The next two calls fail; every call after that succeeds.
api.SetErrors(thirdparty.ErrServiceUnavailable, &thirdparty.TooManyRequestsError{RetryAfter: time.Second})

api.Post(ctx, "a") // error
api.Post(ctx, "b") // error
api.Post(ctx, "c") // nil
api.Posted()       // [c] – only successful posts are recorded
```

```sh
cd go
go build ./...
go test -race ./...
go vet ./...
go run .
```

## Getting Started

Open the repository in the dev container (see [`.devcontainer/README.md`](.devcontainer/README.md)), then:

```sh
claude
```
