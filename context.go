package hoff

import (
	"github.com/google/go-cmp/cmp"
)

type Context struct {
	data contextData
}

func (x Context) Equal(y Context) bool {
	return cmp.Equal(x.data, y.data)
}

type contextData map[string]interface{}

func (c *Context) Store(key string, value interface{}) {
	c.data[key] = value
}

func (c *Context) Delete(key string) {
	delete(c.data, key)
}

func (c *Context) Read(key string) (interface{}, bool) {
	value, ok := c.data[key]
	return value, ok
}

func (c *Context) HaveKey(key string) bool {
	_, ok := c.data[key]
	return ok
}

func NewContext() *Context {
	return setupContext(make(contextData))
}

func setupContext(data contextData) *Context {
	if data == nil {
		return NewContext()
	}
	return &Context{
		data: data,
	}
}
