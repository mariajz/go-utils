package inject

import (
	"errors"
	"fmt"
	"reflect"
)

type Container struct {
	instances map[string]reflect.Value
}

func NewContainer() *Container {
	return &Container{
		instances: make(map[string]reflect.Value),
	}
}

func (c *Container) Provide(key string, instance interface{}) {
	val := reflect.ValueOf(instance)
	if val.Kind() != reflect.Ptr {
		panic(errors.New("instance must be a pointer"))
	}
	c.instances[key] = val
	return
}

func (c *Container) Register(target interface{}) {
	targetValue := reflect.ValueOf(target)
	if targetValue.Kind() != reflect.Ptr || targetValue.Elem().Kind() != reflect.Struct {
		panic(errors.New("target must be a pointer to a struct"))
	}
	targetType := targetValue.Elem().Type()

	for i := 0; i < targetType.NumField(); i++ {
		field := targetType.Field(i)
		tag := field.Tag.Get("inject")
		if tag != "" {
			var filteredKey string
			for k := range c.instances {
				if k == tag {
					filteredKey = k
					break
				}
			}

			if filteredKey != "" {
				instanceValue, ok := c.instances[filteredKey]
				if !ok {
					panic(fmt.Errorf("no instance found for %s", tag))
				}
				targetValue.Elem().Field(i).Set(instanceValue)

			} else {
				panic(fmt.Errorf("no instance found for %s", tag))
			}
		}
	}
}
