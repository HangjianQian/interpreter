package core

import (
	"fmt"
)

type Env struct {
	enclosing *Env
	values    map[string]interface{}
}

func NewEnv(enclosing *Env) *Env {
	return &Env{
		enclosing: enclosing,
		values:    make(map[string]interface{}),
	}
}

func (e *Env) assign(t Token, v interface{}) {
	if _, ok := e.values[t.lexeme]; ok {
		e.values[t.lexeme] = v
		return
	}

	if e.enclosing != nil {
		e.enclosing.assign(t, v)
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
	if e.enclosing != nil {
		return e.enclosing.get(t)
	}
	panic(fmt.Sprintf("get, map key not exist: %s", t.lexeme))
}

func (e *Env) getAt(distance int, name string) interface{} {
	return e.ancestor(distance).values[name]
}

func (e *Env) assignAt(distance int, name string, v interface{}) {
	e.ancestor(distance).values[name] = v
}

func (e *Env) ancestor(distance int) *Env {
	var en *Env = e
	for i := 0; i < distance; i++ {
		en = en.enclosing
	}
	return en
}
