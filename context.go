package hoff

import (
	"github.com/google/go-cmp/cmp"
)

// Context hold data during an Computation
type Context struct {
	Data map[string]interface{}
}

// NewContextWithoutData generate a new empty Context
func NewContextWithoutData() *Context {
	return &Context{
		Data: make(map[string]interface{}),
	}
}

// NewContext generate a new Context with data
func NewContext(data map[string]interface{}) *Context {
	return &Context{
		Data: data,
	}
}

// Equal validate the two Context are equals
func (c Context) Equal(o Context) bool {
	return cmp.Equal(c.Data, o.Data)
}

// Store add a key and its value to the context
func (c *Context) Store(key string, value interface{}) {
	c.Data[key] = value
}

// Delete remove a value in the context by its key
func (c *Context) Delete(key string) {
	delete(c.Data, key)
}

// Read get a value in the context by its key
func (c *Context) Read(key string) (interface{}, bool) {
	value, ok := c.Data[key]
	return value, ok
}

// HaveKey validate that a key is in the context
func (c *Context) HaveKey(key string) bool {
	_, ok := c.Data[key]
	return ok
}
