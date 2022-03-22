package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/assertive-lang/asserlang/Asserlang_go/lexer"

	"github.com/assertive-lang/asserlang/Asserlang_go/parser"
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
	p := parser.New(l)
	program := p.ParseProgram()

	if program == nil {
		panic(fmt.Errorf("ParseProgram() returned nil"))
	}
	errors := p.Errors()
	if len(errors) > 0 {
		for _, msg := range errors {
			fmt.Printf(msg)
		}
	}

	for i, s := range program.Statements {
		fmt.Printf("%d | %s\n", i+1, s)
	}
}
