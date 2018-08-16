package flowengine

import (
	"testing"
)

func isKeyPresent(key string) func(*FlowContext) bool {
	return func(c *FlowContext) bool {
		_, err := c.Read(key)
		return err == nil
	}
}

func Test_DecisionNode(t *testing.T) {
	tc := []FlowNodeTestCase{
		{
			"True",
			contextData{
				"key": "value",
			},
			NewDecisionNode(
				isKeyPresent("key"),
			),
			contextData{},
		},
		{
			"False",
			contextData{},
			NewDecisionNode(
				isKeyPresent("key"),
			),
			contextData{},
		},
	}
	RunTestOnFlowNode(t, tc)
}
