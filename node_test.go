package flowengine

import (
	"testing"
)

type SomeNode struct{}

func (n *SomeNode) Run(c *FlowContext) RunState {
	c.Store("message", "SomeNode is passing")
	return RunStatePass()
}
func (n *SomeNode) AvailableBranches() []NodeBranch {
	return AvailablesBranches()
}

func Test_SomeNode(t *testing.T) {
	tc := []FlowNodeTestCase{
		{
			"SomeNode is passing",
			contextData{},
			&SomeNode{},
			RunStatePass(),
			contextData{
				"message": "SomeNode is passing",
			},
		},
	}
	RunTestOnFlowNode(t, tc)
}
