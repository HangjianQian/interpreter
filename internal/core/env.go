package core

import "fmt"

type Env struct {
	values map[string]interface{}
}

func NewEnv() *Env {
	return &Env{
		values: make(map[string]interface{}),
	}
}

func (e *Env) assign(t Token, v interface{}) {
	if _, ok := e.values[t.lexeme]; ok {
		e.values[t.lexeme] = v
		return
	}
	panic(fmt.Sprintf("assign, map key not exist: %s", t.lexeme))
}

func (e *Env) define(k string, v interface{}) {
	e.values[k] = v
}

func (e *Env) get(t Token) interface{} {
	if v, ok := e.values[t.lexeme]; ok {
		return v
	}
	panic(fmt.Sprintf("get, map key not exist: %s", t.lexeme))
}
