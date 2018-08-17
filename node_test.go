package flow

import (
	"testing"
)

type SomeNode struct{}

func (n *SomeNode) Run(c *Context) RunState {
	c.Store("message", "SomeNode is passing")
	return RunStatePass()
}
func (n *SomeNode) AvailableBranches() []NodeBranch {
	return nil
}

func Test_SomeNode(t *testing.T) {
	tc := []NodeTestCase{
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
	RunTestOnNode(t, tc)
}
