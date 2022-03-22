package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/assertive-lang/asserlang/Asserlang_go/lexer"
	"github.com/assertive-lang/asserlang/Asserlang_go/token"
)

var FilePath string

func init() {
	flag.Parse()
	FilePath = flag.Args()[0]
}

func main() {
	raw, err := ioutil.ReadFile(FilePath)
	if err != nil {
		panic(err)
	}
	l := lexer.New(string(raw))
	for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
		fmt.Printf("%+v\n", tok)
	}
}
