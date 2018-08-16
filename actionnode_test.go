package flowengine

import (
	"fmt"
	"testing"
)

func Test_ActionNode(t *testing.T) {
	tc := []FlowNodeTestCase{
		{
			"No action",
			contextData{},
			NewActionNode(func(c *FlowContext) {}),
			contextData{},
		},
		{
			"Store in context",
			contextData{},
			NewActionNode(func(c *FlowContext) {
				c.Store("key", "value")
			}),
			contextData{
				"key": "value",
			},
		},
		{
			"Read in context",
			contextData{
				"key": "value",
			},
			NewActionNode(func(c *FlowContext) {
				value, _ := c.Read("key")
				c.Store("stored_key", value)
			}),
			contextData{
				"key":        "value",
				"stored_key": "value",
			},
		},
		{
			"Error on Read in context",
			contextData{},
			NewActionNode(func(c *FlowContext) {
				_, err := c.Read("key")
				c.Store("read_error", err)
			}),
			contextData{
				"read_error": fmt.Errorf("unknown key: key"),
			},
		},
	}
	RunTestOnFlowNode(t, tc)
}
