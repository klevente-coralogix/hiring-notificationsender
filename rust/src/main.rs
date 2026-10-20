use std::time::Duration;

use hiring_notificationsender::thirdparty::{MockApi, ThirdPartyApi};

#[tokio::main]
async fn main() {
    let api = MockApi::new();
    api.set_delay(Duration::from_millis(100));

    let result = api.post("Hello, world!").await;
    println!("post result: {result:?}, posted: {:?}", api.posted());
}
