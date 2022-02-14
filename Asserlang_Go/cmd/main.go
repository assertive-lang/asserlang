package main

import (
	"flag"
	"log"
)

func main() {
	flag.Parse()
	path := flag.Args()[0] // file name
	log.Println(path)
}
