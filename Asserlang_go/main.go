package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/assertive-lang/asserlang/Asserlang_go/eval"
	"github.com/assertive-lang/asserlang/Asserlang_go/lexer"
	"github.com/assertive-lang/asserlang/Asserlang_go/object"
	"github.com/assertive-lang/asserlang/Asserlang_go/repl"

	"github.com/assertive-lang/asserlang/Asserlang_go/parser"
)

var args []string

func init() {
	flag.Parse()
	args = flag.Args()
}

func main() {
	if len(args) > 0 {
		FilePath := args[0]
		raw, err := ioutil.ReadFile(FilePath)
		if err != nil {
			panic(fmt.Errorf("어쩔파일: %s", err))
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
				fmt.Printf("%s\n", msg)
			}

		} else {
			env := object.NewEnvironment()
			eval.Eval(program, env)
			if len(eval.Errors) > 0 {
				for _, err := range eval.Errors {
					println(err.Inspect())
				}
			}
		}

	} else {
		repl.Start()
	}
}
