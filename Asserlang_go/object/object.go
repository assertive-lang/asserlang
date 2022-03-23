package object

import (
	"fmt"
)

const (
	INTEGER_OBJ = "INTEGER"
	ERROR_OBJ   = "ERROR"
	BUILTIN_OBJ = "BUILTIN"
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
}

func bPrint(args ...Object) Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}
	return nil
}
