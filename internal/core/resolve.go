package core

import (
	"fmt"
)

type Resolver struct {
	inter  *Interpreter
	scopes *Scope
}

type Scope []map[string]bool

func (s *Scope) empty() bool {
	return len(*s) == 0
}

func (s *Scope) get(i int) map[string]bool {
	return (*s)[i]
}

func (s *Scope) push(v map[string]bool) {
	*s = append(*s, v)
}

func (s *Scope) pop() (map[string]bool, bool) {
	if s.empty() {
		return nil, false
	}
	t := (*s)[len(*s)-1]
	*s = (*s)[:len(*s)-1]
	return t, true
}

func (s *Scope) peek() map[string]bool {
	return (*s)[len(*s)-1]
}

func (s *Scope) size() int {
	return len(*s)
}

func NewResolver(i *Interpreter) *Resolver {
	r := &Resolver{
		inter:  i,
		scopes: new(Scope),
	}
	return r
}

func (r *Resolver) resolve(s interface{}) {
	switch v := s.(type) {
	case AssignExpr:
		r.resolveAssignExpr(v)
	case BinaryExpr:
		r.resolveBinaryExpr(v)
	case UnaryExpr:
		r.resolveUnaryExpr(v)
	case GroupExpr:
		r.resolveGroupExpr(v)
	case LiteralExpr:
		r.resolveLiteralExpr(v)
	case VarExpr:
		r.resolveVarExpr(v)
	case LogicalExpr:
		r.resolveLogicalExpr(v)
	case CallExpr:
		r.resolveCallExpr(v)
	case VarStmt:
		r.resolveVarStmt(v)
	case ExprStmt:
		r.resolveExprStmt(v)
	case BlockStmt:
		r.resolveBlockStmt(v)
	case IfStmt:
		r.resolveIfStmt(v)
	case WhileStmt:
		r.resolveWhileStmt(v)
	case FuncStmt:
		r.resolveFunctionStmt(v)
	case ReturnStmt:
		r.resolveReturnStmt(v)
	case []Stmt:
		for _, i := range v {
			r.resolve(i)
		}
	default:
		panic(fmt.Sprintf("unimplement resolve, %v", s))
	}
}

func (r *Resolver) resolveBlockStmt(sts BlockStmt) interface{} {
	r.beginScope()
	r.resolve(sts.stmts)
	r.endScope()
	return nil
}

func (r *Resolver) resolveVarStmt(v VarStmt) interface{} {
	r.declare(v.name)
	if v.initializer != nil {
		// expr
		r.resolve(v.initializer)
	}
	r.define(v.name)
	return nil
}

func (r *Resolver) resolveVarExpr(v VarExpr) interface{} {
	if !r.scopes.empty() {
		if res, ok := r.scopes.peek()[v.name.lexeme]; ok && !res {
			panic(fmt.Sprintf("can't read local variable in its own initializer, line: %d, %s", v.name.line, v.name.lexeme))
		}
	}

	r.resolveLocal(v, v.name)
	return nil
}

func (r *Resolver) resolveLocal(ex Expr, n Token) {
	for i := r.scopes.size() - 1; i >= 0; i-- {
		if _, ok := r.scopes.get(i)[n.lexeme]; ok {
			r.inter.resolve(ex, r.scopes.size()-1-i)
			return
		}
	}

}

func (r *Resolver) resolveExprStmt(ex ExprStmt) interface{} {
	r.resolve(ex.expr)
	return nil
}

func (r *Resolver) resolveIfStmt(st IfStmt) interface{} {
	r.resolve(st.condition)
	r.resolve(st.thenBranch)
	if st.elseBranch != nil {
		r.resolve(st.elseBranch)
	}
	return nil
}

func (r *Resolver) resolveReturnStmt(st ReturnStmt) interface{} {
	if st.value != nil {
		r.resolve(st.value)
	}
	return nil
}

func (r *Resolver) resolveWhileStmt(st WhileStmt) interface{} {
	r.resolve(st.condition)
	r.resolve(st.body)
	return nil
}

func (r *Resolver) resolveBinaryExpr(ex BinaryExpr) interface{} {
	r.resolve(ex.left)
	r.resolve(ex.right)
	return nil
}

func (r *Resolver) resolveCallExpr(ex CallExpr) interface{} {
	r.resolve(ex.callee)

	for _, arg := range ex.args {
		r.resolve(arg)
	}
	return nil
}

func (r *Resolver) resolveGroupExpr(ex GroupExpr) interface{} {
	r.resolve(ex.expression)
	return nil
}

func (r *Resolver) resolveLiteralExpr(ex LiteralExpr) interface{} {
	return nil
}

func (r *Resolver) resolveLogicalExpr(ex LogicalExpr) interface{} {
	r.resolve(ex.left)
	r.resolve(ex.right)
	return nil
}

func (r *Resolver) resolveUnaryExpr(ex UnaryExpr) interface{} {
	r.resolve(ex.right)
	return nil
}

func (r *Resolver) resolveFunctionStmt(ex FuncStmt) interface{} {
	r.declare(ex.name)
	r.define(ex.name)

	r.resolveFunction(ex)
	return nil
}

func (r *Resolver) resolveFunction(fc FuncStmt) interface{} {
	r.beginScope()
	for _, arg := range fc.params {
		r.declare(arg)
		r.declare(arg)
	}
	r.resolve(fc.body)
	r.endScope()
	return nil
}

func (r *Resolver) resolveAssignExpr(ex AssignExpr) interface{} {
	r.resolve(ex.value)
	r.resolveLocal(ex, ex.name)
	return nil
}

func (r *Resolver) declare(n Token) {
	if r.scopes.empty() {
		return
	}
	r.scopes.peek()[n.lexeme] = false
}

func (r *Resolver) define(n Token) {
	if r.scopes.empty() {
		return
	}
	r.scopes.peek()[n.lexeme] = true
}

func (r *Resolver) beginScope() {
	r.scopes.push(map[string]bool{})
}

func (r *Resolver) endScope() {
	_, ok := r.scopes.pop()
	if !ok {
		panic("scope is empty")
	}
}
