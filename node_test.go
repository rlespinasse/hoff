package flowengine

import (
	"testing"
)

type SomeNode struct{}

func (n *SomeNode) Run(c *FlowContext) {
	c.Store("message", "SomeNode is running")
}

func Test_SomeNode(t *testing.T) {
	tc := []FlowNodeTestCase{
		{
			"SomeNode is running",
			contextData{},
			&SomeNode{},
			contextData{
				"message": "SomeNode is running",
			},
		},
	}
	RunTestOnFlowNode(t, tc)
}
