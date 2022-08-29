use tiny_http::{Response, Server};

const ADDR: &str = "0.0.0.0";
const PORT: u32 = 5001;
const FORWARDED: &str = "Forwarded";
const FORWARDED_FOR: &str = "X-Forwarded-For";

fn main() {
    println!("starting server on {}:{}", ADDR, PORT);
    let server = Server::http(format!("{}:{}", ADDR, PORT)).unwrap();

    for request in server.incoming_requests() {
        let remote_addr = request.remote_addr();
        let ip = remote_addr.ip();
        let mut real_ip = String::from("");
        for header in request.headers() {
            if header.field.equiv(FORWARDED_FOR) || header.field.equiv(FORWARDED) {
                real_ip = header.value.to_string();
            }
        }

        if real_ip.is_empty() {
            real_ip = ip.to_string();
        }

        let response = Response::from_string(real_ip);
        match request.respond(response) {
            Ok(()) => (),
            Err(err) => println!("failed to send request: {}", err),
        }
    }
}
