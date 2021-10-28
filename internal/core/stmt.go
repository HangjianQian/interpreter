package core

type Stmt interface{}

type ExprStmt struct {
	expr Expr
}

type VarStmt struct {
	name        Token
	initializer Expr
}

func (v VarStmt) evaluate() {
	
}
