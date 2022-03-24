package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/assertive-lang/asserlang/Asserlang_go/eval"
	"github.com/assertive-lang/asserlang/Asserlang_go/lexer"
	"github.com/assertive-lang/asserlang/Asserlang_go/object"
	"github.com/assertive-lang/asserlang/Asserlang_go/parser"
)

func Start() {
	scanner := bufio.NewScanner(os.Stdin)

	env := object.NewEnvironment()
	for {

		fmt.Printf(">> ")
		scanned := scanner.Scan()

		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line + "\n슉슈슉슉")

		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			for _, msg := range p.Errors() {
				fmt.Printf("%s\n", msg)
			}
			continue
		}
		evaluated := eval.Eval(program, env)
		if evaluated != nil {
			io.WriteString(os.Stdout, evaluated.Inspect())
			io.WriteString(os.Stdout, "\n")
		}
	}
}
