package di

import (
	"fmt"
	"reflect"
)

type Container struct {
	instances map[reflect.Type]interface{}
}

func NewContainer() *Container {
	return &Container{
		instances: make(map[reflect.Type]interface{}),
	}
}

func (c *Container) Register(instance interface{}) {
	t := reflect.TypeOf(instance)
	if _, exists := c.instances[t]; exists {
		panic(fmt.Sprintf("instance already registered for type %s", t))
	}
	c.instances[t] = instance
}

func (c *Container) Resolve(target interface{}) {
	ptr := reflect.ValueOf(target)
	if ptr.Kind() != reflect.Ptr || ptr.IsNil() {
		panic("target must be a non-nil pointer")
	}

	value := ptr.Elem()
	if value.Kind() != reflect.Interface && value.Kind() != reflect.Ptr {
		panic("target must be a pointer to an interface or pointer")
	}

	instance, exists := c.instances[value.Type()]
	if !exists {
		panic(fmt.Sprintf("no instance found for type %s", value.Type()))
	}

	ptr.Elem().Set(reflect.ValueOf(instance))
}
