package core

type Expr interface{}

type BinaryExpr struct {
	left     Expr
	right    Expr
	operator Token
}

type UnaryExpr struct {
	operator Token
	right    Expr
}

type GroupExpr struct {
	expression Expr
}

type LiteralExpr struct {
	obj interface{}
}
