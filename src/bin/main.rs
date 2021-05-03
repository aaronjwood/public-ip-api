use public_ip_api::ThreadPool;
use std::io;
use std::io::prelude::*;
use std::io::BufReader;
use std::net::TcpListener;
use std::net::TcpStream;

const ADDR: &str = "0.0.0.0";
const PORT: u32 = 5001;
const FORWARDED_HEADER_KEY: &str = "x-forwarded-for";
const RESPONSE_OK: &str = "HTTP/1.1 200 OK\r\n\r\n";
const RESPONSE_INTERNAL_SERVER_ERROR: &str = "HTTP/1.1 200 OK\r\n\r\n";

fn main() -> io::Result<()> {
    TcpListener::bind(format!("{}:{}", ADDR, PORT)).and_then(|listener| {
        let pool = ThreadPool::new(32);
        for stream in listener.incoming() {
            let tcp_stream = stream?;
            pool.execute(|| {
                if let Err(err) = handle_connection(tcp_stream) {
                    println!("{}", err);
                }
            });
        }

        Ok(listener)
    })?;

    Ok(())
}

fn get_header<'a>(key: &str, line: &'a str) -> Result<&'a str, ()> {
    let not_found = Err(());
    let end = line.chars().position(|c| c == ':');
    match end {
        Some(idx) => {
            let header = &line[..idx];
            if key == header.to_lowercase() {
                return Ok(&line[idx + 2..]);
            }

            return not_found;
        }
        None => return not_found,
    }
}

fn handle_connection(stream: TcpStream) -> Result<(), io::Error> {
    let reader = BufReader::new(&stream);
    let mut ip: String;
    let addr = stream.local_addr();
    match addr {
        Ok(sock_addr) => ip = sock_addr.ip().to_string(),
        Err(_) => {
            send_response(stream, RESPONSE_INTERNAL_SERVER_ERROR)?;
            return Ok(());
        }
    }

    let req = parse_request(reader);
    match req {
        Ok(x) => {
            if x != "" {
                ip = x.to_string();
            }
        }
        Err(_) => {
            send_response(stream, RESPONSE_INTERNAL_SERVER_ERROR)?;
            return Ok(());
        }
    }

    let res = format!("{}{}", RESPONSE_OK, ip);
    send_response(stream, &res)?;
    return Ok(());
}

fn parse_request(reader: BufReader<&TcpStream>) -> Result<String, ()> {
    for chunk in reader.lines() {
        match chunk {
            Ok(line) => {
                if line == "" {
                    break;
                }

                let header = get_header(FORWARDED_HEADER_KEY, &line);
                match header {
                    Ok(header) => {
                        let mut original_ip_parts = header.split(", ");
                        let original_ip = original_ip_parts.nth(0);
                        match original_ip {
                            Some(x) => {
                                return Ok(x.to_string());
                            }
                            None => break,
                        }
                    }
                    Err(_) => continue,
                }
            }
            Err(_) => return Err(()),
        };
    }

    return Ok("".to_string());
}

fn send_response(mut stream: TcpStream, response: &str) -> Result<(), io::Error> {
    stream.write(response.as_bytes())?;
    stream.flush()?;
    Ok(())
}
