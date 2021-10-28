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

func (e *Env) define(k string, v interface{}) {
	e.values[k] = v
}

func (e *Env) get(t Token) interface{} {
	if v, ok := e.values[t.lexeme]; ok {
		return v
	}
	panic(fmt.Sprintf("map key not exist: ", t.lexeme))
}
