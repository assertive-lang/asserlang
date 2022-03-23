package eval

import (
	"fmt"

	"github.com/assertive-lang/asserlang/Asserlang_go/ast"
	"github.com/assertive-lang/asserlang/Asserlang_go/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {

	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		env.Set(node.Name.Value, val)

		return val

	case *ast.InfixIntegerExpression:
		leftObj := Eval(node.Left, env)
		return &object.Integer{Value: leftObj.(*object.Integer).Value + node.Right}

	case *ast.TUExpression:
		leftObj := Eval(node.Left, env)
		rightObj := Eval(node.Right, env)
		return &object.Integer{Value: leftObj.(*object.Integer).Value * rightObj.(*object.Integer).Value}

	case *ast.IntegerLiteral:
		return &object.Integer{Value: node.Value}

	case *ast.Identifier:
		if val, ok := env.Get(node.Value); ok {
			return val
		}

		if builtinFn, ok := builtinFunctions[node.Value]; ok {
			return builtinFn
		}

		return newError("어쩔변수")

	}
	return nil
}

func evalStatements(stmts []ast.Statement, env *object.Environment) object.Object {
	var result object.Object
	for _, statement := range stmts {
		result = Eval(statement, env)
	}
	return result
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

var builtinFunctions = map[string]*object.Builtin{
	"print": object.GetBuiltinByName("print"),
}
