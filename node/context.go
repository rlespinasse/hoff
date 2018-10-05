package node

import (
	"github.com/google/go-cmp/cmp"
)

// Context hold data during an Computation
type Context struct {
	data map[string]interface{}
}

// NewWithoutData generate a new empty Context
func NewWithoutData() *Context {
	return &Context{
		data: make(map[string]interface{}),
	}
}

// New generate a new Context with data
func New(data map[string]interface{}) *Context {
	return &Context{
		data: data,
	}
}

// Equal validate the two Context are equals
func (c Context) Equal(o Context) bool {
	return cmp.Equal(c.data, o.data)
}

// Store add a key and its value to the context
func (c *Context) Store(key string, value interface{}) {
	c.data[key] = value
}

// Delete remove a value in the context by its key
func (c *Context) Delete(key string) {
	delete(c.data, key)
}

// Read get a value in the context by its key
func (c *Context) Read(key string) (interface{}, bool) {
	value, ok := c.data[key]
	return value, ok
}

// HaveKey validate that a key is in the context
func (c *Context) HaveKey(key string) bool {
	_, ok := c.data[key]
	return ok
}
