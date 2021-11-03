package core

import "fmt"

type Interpreter struct {
	env *Env
}

func NewInterpreter() *Interpreter {
	i := &Interpreter{
		env: NewEnv(nil),
	}

	i.env.define("clock", clockFunc{})
	i.env.define("println", printlnFunc{})

	return i
}

func (i *Interpreter) interpret(s interface{}) interface{} {
	// TODO: fix
	switch v := s.(type) {
	case AssignExpr:
		return i.evaluateAssignExpr(v)
	case BinaryExpr:
		return i.evaluateBinaryExpr(v)
	case UnaryExpr:
		return i.evaluateUnaryExpr(v)
	case GroupExpr:
		return v.evaluate()
	case LiteralExpr:
		return v.evaluate()
	case VarExpr:
		return i.evaluateVarExpr(v)
	case LogicalExpr:
		return i.evaluateLogicalStmt(v)
	case CallExpr:
		return i.evaluateCallExpr(v)
	case VarStmt:
		return i.evaluateVarStmt(v)
	case ExprStmt:
		return i.evaluateExprStmt(v)
	case BlockStmt:
		return i.evaluateBlockStmt(v, NewEnv(i.env))
	case IfStmt:
		return i.evaluateIfStmt(v)
	case WhileStmt:
		return i.evaluateWhileStmt(v)
	case FuncStmt:
		return i.evaluateFuncStmt(v)
	case ReturnStmt:
		return i.evaluateReturnStmt(v)
	}
	return nil
}

func (i *Interpreter) evaluateUnaryExpr(u UnaryExpr) interface{} {
	switch u.operator.kind {
	case MINUS:
		return -1 * i.interpret(u.right).(float64)
	case BANG:
		return !i.interpret(u.right).(bool)
	}
	return nil
}

func (i *Interpreter) evaluateAssignExpr(a AssignExpr) interface{} {
	value := i.interpret(a.value)
	i.env.assign(a.name, value)
	return value
}

func (i *Interpreter) evaluateBinaryExpr(b BinaryExpr) interface{} {
	left := i.interpret(b.left)
	right := i.interpret(b.right)

	switch b.operator.kind {
	case MINUS:
		return left.(float64) - right.(float64)
	case STAR:
		return left.(float64) * right.(float64)
	case SLASH:
		return left.(float64) / right.(float64)
	case PLUS:
		switch left.(type) {
		case float64:
			return left.(float64) + right.(float64)
		case string:
			return left.(string) + right.(string)
		}
	case BANG_EQUAL:
		return !isEqual(left, right)
	case EQUAL_EQUAL:
		return isEqual(left, right)
	case GREATER:
		return left.(float64) > right.(float64)
	case GREATER_EQUAL:
		return left.(float64) >= right.(float64)
	case LESS:
		return left.(float64) < right.(float64)
	case LESS_EQUAL:
		return left.(float64) <= right.(float64)
	}

	// unreachable
	return nil
}

func (i *Interpreter) evaluateVarExpr(v VarExpr) interface{} {
	return i.env.get(v.name)
}

func (i *Interpreter) evaluateVarStmt(v VarStmt) interface{} {
	var obj interface{}
	if v.initializer != nil {
		obj = i.interpret(v.initializer)
	}
	i.env.define(v.name.lexeme, obj)
	return nil
}

func (i *Interpreter) evaluateExprStmt(v ExprStmt) interface{} {
	i.interpret(v.expr)
	return nil
}

func (i *Interpreter) evaluateBlockStmt(v BlockStmt, e *Env) interface{} {
	previousEnv := i.env
	i.env = e
	var (
		err      ReturnErr
		returnOK bool
	)
	for _, stmt := range v.stmts {
		t := i.interpret(stmt)
		if t != nil {
			if err, returnOK = t.(ReturnErr); returnOK {
				break
			}
		}
	}
	i.env = previousEnv
	if returnOK {
		return err.value
	}
	return nil
}

func (i *Interpreter) evaluateIfStmt(v IfStmt) interface{} {
	// TODO: support more condition check, eg: string, float...
	if i.interpret(v.condition).(bool) {
		i.interpret(v.thenBranch)
	} else if v.elseBranch != nil {
		i.interpret(v.elseBranch)
	}
	return nil
}

func (i *Interpreter) evaluateLogicalStmt(v LogicalExpr) interface{} {
	left := i.interpret(v.left).(bool)
	if v.operator.kind == OR {
		if left {
			return true
		}
	} else if !left {
		return false
	}
	return i.interpret(v.right).(bool)
}

func (i *Interpreter) evaluateWhileStmt(v WhileStmt) interface{} {
	for i.interpret(v.condition).(bool) {
		i.interpret(v.body)
	}
	return nil
}

func (i *Interpreter) evaluateCallExpr(v CallExpr) interface{} {
	callee := i.interpret(v.callee)

	var args []interface{}
	for _, a := range v.args {
		args = append(args, i.interpret(a))
	}
	if fn, ok := callee.(Callalble); ok {
		if fn.arity() != len(args) {
			panic(fmt.Sprintf("args num not match, require %d, got %d, line %d", fn.arity(), len(args), v.paren.line))
		}
		return fn.call(i, args)
	} else {
		panic(fmt.Sprintf("invalid fun call at line %d", v.paren.line))
	}
}

func (i *Interpreter) evaluateFuncStmt(v FuncStmt) interface{} {
	fn := v
	i.env.define(fn.name.lexeme, fn)
	return nil
}

func (i *Interpreter) evaluateReturnStmt(v ReturnStmt) interface{} {
	if v.value != nil {
		return ReturnErr{value: i.interpret(v.value)}
	}
	return ReturnErr{value: nil}
}
