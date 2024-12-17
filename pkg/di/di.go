package di

import (
	"errors"

	"go.uber.org/dig"
)

type (
	// Constructor uber dig dependency constructor
	Constructor struct {
		Name string
		Fn   interface{}
	}
)

var container *dig.Container

// Provide teaches the container how to build values of one or more types and expresses their dependencies.
func Provide(constructor interface{}, opts ...dig.ProvideOption) error {
	if container == nil {
		container = dig.New()
	}

	return container.Provide(constructor, opts...)
}

// Invoke runs the given function after instantiating its dependencies.
func Invoke(fn interface{}, opts ...dig.InvokeOption) error {
	if container == nil {
		return errors.New("no constructor provided")
	}
	return container.Invoke(fn, opts...)
}
