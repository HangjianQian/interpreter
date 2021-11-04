package core

type Stmt interface{}

type ExprStmt struct {
	expr Expr
}

type VarStmt struct {
	name        Token
	initializer Expr
}

type BlockStmt struct {
	stmts []Stmt
}

type IfStmt struct {
	condition  Expr
	thenBranch Stmt
	elseBranch Stmt
}

type WhileStmt struct {
	condition Expr
	body      Stmt
}

type FuncStmt struct {
	name    Token
	params  []Token
	body    []Stmt
	closure *Env
}

type ReturnStmt struct {
	keyword Token
	value   Expr
}
