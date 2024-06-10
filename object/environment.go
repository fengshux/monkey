package object

import (
	"fmt"
	"strings"
)

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}

	return obj, ok
}

func (e *Environment) Set(name string, obj Object) Object {
	e.store[name] = obj
	return obj
}

func (e *Environment) String() string {
	var out strings.Builder
	out.WriteString("{")

	kvs := make([]string, 0, len(e.store)+1)
	for k, v := range e.store {
		kvs = append(kvs, fmt.Sprintf("%s: %s", k, v.Inspect()))
	}

	if e.outer != nil {
		kvs = append(kvs, fmt.Sprintf("%s: %s", "outer", e.outer.String()))
	}

	out.WriteString(strings.Join(kvs, ", "))

	return out.String()
}
