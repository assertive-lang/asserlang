package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/assertive-lang/asserlang/Asserlang_go/ast"
)

const (
	INTEGER_OBJ     = "INTEGER"
	ERROR_OBJ       = "ERROR"
	BUILTIN_OBJ     = "BUILTIN"
	FUNCTION_OBJ    = "FUNCTION"
	RETURNVALUE_OBJ = "RETURN_VALUE"
	NULL_OBJ        = "NULL"
)

type ObjectType string

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType { return INTEGER_OBJ }

type Error struct {
	Message string
}

// Null type is an empty struct
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }

func (n *Null) Inspect() string { return "null" }

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return e.Message }

type BuiltinFunction func(args ...Object) Object

type Builtin struct {
	Fn BuiltinFunction
}

func (n *Builtin) Type() ObjectType { return BUILTIN_OBJ }

func (n *Builtin) Inspect() string { return "builtin function" }

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}

	return nil
}

func newError(msgWithFormatVerbs string, values ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(msgWithFormatVerbs, values...)}
}

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{"ㅇㅉ", &Builtin{Fn: bPrint}},
	{"ㅌㅂ", &Builtin{Fn: bInput}},
}

func bPrint(args ...Object) Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}
	return nil
}

func bInput(args ...Object) Object {
	var tmp string
	fmt.Scanln(&tmp)
	return nil
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }

func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("funclmao")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n")

	return out.String()
}

type ReturnValue struct {
	Value Object
}

func (rv *ReturnValue) Type() ObjectType { return RETURNVALUE_OBJ }

func (rv *ReturnValue) Inspect() string { return rv.Value.Inspect() }
