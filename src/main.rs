use actix_web::{web, App, HttpRequest, HttpServer};

const ADDR: &str = "0.0.0.0";
const PORT: u32 = 5001;

async fn handler(req: HttpRequest) -> String {
    match req.connection_info().realip_remote_addr() {
        Some(ip) => match ip.split(":").next() {
            Some(part) => part.to_string(),
            None => {
                let err = "failed to parse IP";
                println!("{}", err);
                err.to_owned()
            }
        },
        None => {
            let err = "failed to get remote address";
            println!("{}", err);
            err.to_owned()
        }
    }
}

#[actix_web::main]
async fn main() -> std::io::Result<()> {
    println!("starting server on {}:{}", ADDR, PORT);
    HttpServer::new(|| App::new().route("/", web::get().to(handler)))
        .bind(format!("{}:{}", ADDR, PORT))?
        .run()
        .await
}
