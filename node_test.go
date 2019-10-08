package hoff

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/computestate"
)

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
