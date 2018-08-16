package flowengine

import "fmt"

type FlowContext struct {
	data contextData
}

type contextData map[string]interface{}
type DataValue interface{}

func (c *FlowContext) Store(key string, value interface{}) {
	c.data[key] = value
}

func (c *FlowContext) Read(key string) (interface{}, error) {
	value, ok := c.data[key]
	if ok {
		return value, nil
	}
	return nil, fmt.Errorf("unknown key: %s", key)
}

func NewFlowContext() *FlowContext {
	return setupFlowContext(make(contextData))
}

func setupFlowContext(data contextData) *FlowContext {
	return &FlowContext{
		data: data,
	}
}
