package eval

import (
	"fmt"

	"github.com/assertive-lang/asserlang/Asserlang_go/ast"
	"github.com/assertive-lang/asserlang/Asserlang_go/object"
)

func Eval(node ast.Node, env *object.Environment) object.Object {
	// fmt.Printf("%T\n", node)

	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements, env)

	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)

	case *ast.LetStatement:
		val := Eval(node.Value, env)
		env.Set(node.Name.Value, val)

		return val

	case *ast.BlockStatement:
		return evalBlockStmt(node, env)

	case *ast.FunctionLiteral:
		params := node.Parameters
		body := node.Body
		val := &object.Function{
			Parameters: params,
			Body:       body,
			Env:        env,
		}
		env.Set(node.Name, val)
		return val

	case *ast.CallExpression:
		fn := Eval(node.Function, env)
		if isError(fn) {
			return fn
		}
		args := evalExprs(node.Arguments, env)
		if len(args) == 1 && isError(args[0]) {
			return args[0]
		}
		return applyFunction(fn, args, node.Token.Line)

	case *ast.InfixIntegerExpression:
		leftObj := Eval(node.Left, env)

		switch left := leftObj.(type) {
		case *object.Integer:
			return &object.Integer{Value: left.Value + node.Right}
		case *object.Builtin:

			return applyFunction(left, []object.Object{&object.Integer{Value: node.Right}}, node.Token.Line)

		}

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

		return newError(fmt.Sprintf("어쩔변수 at line %d: %s", node.Token.Line, node.Value))

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

func evalBlockStmt(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object

	for _, stmt := range block.Statements {
		result = Eval(stmt, env)

		if result != nil {
			rt := result.Type()
			if rt == object.RETURNVALUE_OBJ || rt == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

var builtinFunctions = map[string]*object.Builtin{
	"ㅇㅉ": object.GetBuiltinByName("ㅇㅉ"),
}

func applyFunction(function object.Object, args []object.Object, line int) object.Object {
	switch fn := function.(type) {
	case *object.Function:
		extendedEnv := extendFunctionEnv(fn, args)
		evaluated := Eval(fn.Body, extendedEnv)
		return unwrapReturnValue(evaluated)
	case *object.Builtin:
		if result := fn.Fn(args...); result != nil {
			return result
		}
		return nil
	default:
		return newError("Line %d: Not a function: %s", line, function.Type())
	}
}

func extendFunctionEnv(fn *object.Function, args []object.Object) *object.Environment {
	env := object.NewEnclosedEnvironment(fn.Env)

	for i, param := range fn.Parameters {
		env.Set(param.Value, args[i])
	}
	return env
}

func unwrapReturnValue(obj object.Object) object.Object {
	if returnValue, ok := obj.(*object.ReturnValue); ok {
		return returnValue.Value
	}
	return obj
}

func evalExprs(exprs []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object

	for _, expr := range exprs {
		evaluated := Eval(expr, env)
		if isError(evaluated) {
			return []object.Object{evaluated}
		}
		result = append(result, evaluated)
	}

	return result
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
