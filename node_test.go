package flow

import (
	"testing"
)

type SomeNode struct{}

func (n *SomeNode) Compute(c *Context) ComputeState {
	c.Store("message", "SomeNode is passing")
	return ComputeStatePass()
}
func (n *SomeNode) AvailableBranches() []string {
	return nil
}

func Test_SomeNode(t *testing.T) {
	tc := []NodeTestCase{
		{
			"SomeNode is passing",
			contextData{},
			&SomeNode{},
			ComputeStatePass(),
			contextData{
				"message": "SomeNode is passing",
			},
		},
	}
	ComputeTestOnNode(t, tc)
}
