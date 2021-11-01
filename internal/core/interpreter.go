package core

type Interpreter struct {
	env *Env
}

func NewInterpreter() *Interpreter {
	return &Interpreter{
		env: NewEnv(),
	}
}

func (i *Interpreter) interpret(s interface{}) interface{} {
	// TODO: fix
	switch v := s.(type) {
	case AssignExpr:
		return i.evaluateAssignExpr(v)
	case BinaryExpr:
		return i.evaluateBinaryExpr(v)
	case UnaryExpr:
		return v.evaluate()
	case GroupExpr:
		return v.evaluate()
	case LiteralExpr:
		return v.evaluate()
	case VarExpr:
		return i.evaluateVarExpr(v)
	case VarStmt:
		return i.evaluateVarStmt(v)
	case ExprStmt:
		return i.evaluateExprStmt(v)
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
