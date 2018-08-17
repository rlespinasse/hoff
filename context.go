package flow

import "fmt"

type Context struct {
	data contextData
}

type contextData map[string]interface{}

func (c *Context) Store(key string, value interface{}) {
	c.data[key] = value
}

func (c *Context) Read(key string) (interface{}, error) {
	value, ok := c.data[key]
	if ok {
		return value, nil
	}
	return nil, fmt.Errorf("unknown key: %s", key)
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
