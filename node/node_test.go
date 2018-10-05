package node

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/computestate"
)

type SomeNode struct{}

func (n *SomeNode) Compute(c *Context) computestate.ComputeState {
	c.Store("message", "SomeNode is passing")
	return computestate.Continue()
}

func (n *SomeNode) DecideCapability() bool {
	return false
}

type AnotherNode struct{}

func (n *AnotherNode) Compute(c *Context) computestate.ComputeState {
	c.Store("message", "AnotherNode is passing")
	return computestate.Continue()
}

func (n *AnotherNode) DecideCapability() bool {
	return false
}

func Test_SomeNode(t *testing.T) {
	tc := []NodeTestCase{
		{
			name:                 "Should SomeNode Pass and Store a message",
			givenNode:            &SomeNode{},
			expectedComputeState: computestate.Continue(),
			expectedContextData:  map[string]interface{}{"message": "SomeNode is passing"},
		},
	}
	RunTestOnNode(t, tc)
}

func Test_NodeComparator_Equal(t *testing.T) {
	givenNode := &SomeNode{}
	givenAnotherNode := givenNode

	if !cmp.Equal(givenNode, givenAnotherNode, NodeComparator) {
		t.Errorf("node: %+v and anotherNode: %+v must be equals", givenNode, givenAnotherNode)
	}
}

func Test_NodeComparator_NotEqual(t *testing.T) {
	givenNode := &SomeNode{}
	givenAnotherNode := &AnotherNode{}

	if cmp.Equal(givenNode, givenAnotherNode, NodeComparator) {
		t.Errorf("node: %+v and anotherNode: %+v must not be equals", givenNode, givenAnotherNode)
	}
}
