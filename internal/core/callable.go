package core

import (
	"fmt"
	"time"
)

type Callalble interface {
	arity() int
	call(i *Interpreter, args []interface{}) interface{}
}

type clockFunc struct{}

func (c clockFunc) arity() int {
	return 0
}

func (c clockFunc) call(i *Interpreter, args []interface{}) interface{} {
	return float64(time.Now().UTC().Second())
}

type printlnFunc struct{}

func (p printlnFunc) arity() int {
	return 1
}

func (p printlnFunc) call(i *Interpreter, args []interface{}) interface{} {
	fmt.Println(args)
	return nil
}

func (f FuncStmt) arity() int {
	return len(f.params)
}

func (f FuncStmt) call(i *Interpreter, args []interface{}) interface{} {
	env := NewEnv(i.env)
	for idx, _ := range f.params {
		env.define(f.params[idx].lexeme, args[idx])
	}

	i.evaluateBlockStmt(BlockStmt{f.body}, env)
	return nil
}
