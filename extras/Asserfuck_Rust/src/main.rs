use std::fs;
use std::env;
use std::collections::HashMap;
use std::path::Path;
use std::vec::Vec;
use std::fs::File;
use std::io::prelude::*;

fn main() {
    let keywords: HashMap<&str, &str> = HashMap::from([
        ("ㅋ", "+"),
        ("ㅎ", "-"),
        ("어", ">"),
        ("쩔", "<"),
        ("저", "."),
        ("쪌", ","),
        ("티", "["),
        ("비", "]")
    ]);
    let filename = env::args().nth(1).expect("Nothing given");
    let new_file_name = format!("{filename}.bf");
    let path = Path::new(&new_file_name);
    let display = path.display();
    
    let src = fs::read_to_string(filename)
    .expect("Something went wrong with reading the file");
    let codes = src.split("");
    let mut compiled: Vec<&str> = Vec::new();
    for code in codes {
        if keywords.contains_key(code) {
            compiled.push(keywords[code]);
        }
    }
    let mut file = match File::create(&path) {
        Err(why) => panic!("Couldn't create {}: {}", display, why),
        Ok(file) => file,
    };

    // Write the `LOREM_IPSUM` string to `file`, returns `io::Result<()>`
    match file.write_all(compiled.join("").as_bytes()) {
        Err(why) => panic!("Couldn't write to {}: {}", display, why),
        Ok(_) => println!("Successfully wrote to {}", display),
    }
}
