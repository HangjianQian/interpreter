package core

type Expr interface{}

type BinaryExpr struct {
	left     Expr
	right    Expr
	operator Token
}

func (b BinaryExpr) evaluate() interface{} {
	left := evaluate(b.left)
	right := evaluate(b.right)

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

type UnaryExpr struct {
	operator Token
	right    Expr
}

func (u UnaryExpr) evaluate() interface{} {
	switch u.operator.kind {
	case MINUS:
		return -1 * u.operator.literal.(float64)
	case BANG:
		return !u.operator.literal.(bool)
	}
	return nil
}

type GroupExpr struct {
	expression Expr
}

func (g GroupExpr) evaluate() interface{} {
	return evaluate(g.expression)
}

type LiteralExpr struct {
	obj interface{}
}

func (l LiteralExpr) evaluate() interface{} {
	return l.obj
}

type VarExpr struct {
	name Token
}

type AssignExpr struct {
	name  Token
	value Expr
}

// deprecate
func evaluate(e Expr) interface{} {
	switch v := e.(type) {
	case BinaryExpr:
		return v.evaluate()
	case UnaryExpr:
		return v.evaluate()
	case GroupExpr:
		return v.evaluate()
	case LiteralExpr:
		return v.evaluate()
	case VarExpr:
	}
	return nil
}

func isEqual(a, b interface{}) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil {
		return false
	}

	switch a.(type) {
	case float64:
		// TODO: precition
		return a.(float64) == b.(float64)
	case string:
		return a.(string) == b.(string)
	}
	return false
}
