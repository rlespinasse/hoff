package flow

import (
	"errors"
	"testing"
)

func Test_ActionNode_Compute(t *testing.T) {
	tc := []NodeTestCase{
		{
			name:                 "Should Pass",
			givenNode:            NewActionNode(func(*Context) (bool, error) { return true, nil }),
			expectedComputeState: ComputeStatePass(),
		},
		{
			name:                 "Should Stop",
			givenNode:            NewActionNode(func(*Context) (bool, error) { return false, nil }),
			expectedComputeState: ComputeStateStop(),
		},
		{
			name:                 "Should Fail",
			givenNode:            NewActionNode(func(*Context) (bool, error) { return false, errors.New("error") }),
			expectedComputeState: ComputeStateFail(errors.New("error")),
		},
	}
	RunTestOnNode(t, tc)
}

func Test_ActionNode_AvailableBranches(t *testing.T) {
	node := NewActionNode(nil)
	branches := node.AvailableBranches()

	if branches != nil {
		t.Errorf("got: %+v, want: nil", branches)
	}
}
