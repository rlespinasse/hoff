package namingishard

import (
	"testing"
)

type SomeNode struct{}

func (n *SomeNode) Compute(c *Context) ComputeState {
	c.Store("message", "SomeNode is passing")
	return ComputeStateContinue()
}

func (n *SomeNode) decideCapability() bool {
	return false
}

func Test_SomeNode(t *testing.T) {
	tc := []NodeTestCase{
		{
			name:                 "Should SomeNode Pass and Store a message",
			givenNode:            &SomeNode{},
			expectedComputeState: ComputeStateContinue(),
			expectedContextData:  contextData{"message": "SomeNode is passing"},
		},
	}
	RunTestOnNode(t, tc)
}
