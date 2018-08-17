package flow

import (
	"errors"
	"testing"
)

func Test_ActionNode_Compute(t *testing.T) {
	tc := []NodeTestCase{
		{
			name:                 "Pass",
			givenContextData:     contextData{},
			givenNode:            NewActionNode(func(*Context) (bool, error) { return true, nil }),
			expectedComputeState: ComputeStatePass(),
			expectedContextData:  contextData{},
		},
		{
			name:                 "Stop",
			givenContextData:     contextData{},
			givenNode:            NewActionNode(func(*Context) (bool, error) { return false, nil }),
			expectedComputeState: ComputeStateStop(),
			expectedContextData:  contextData{},
		},
		{
			name:                 "Fail",
			givenContextData:     contextData{},
			givenNode:            NewActionNode(func(*Context) (bool, error) { return false, errors.New("error") }),
			expectedComputeState: ComputeStateFail(errors.New("error")),
			expectedContextData:  contextData{},
		},
	}
	ComputeTestOnNode(t, tc)
}

func Test_ActionNode_AvailableBranches(t *testing.T) {
	node := NewActionNode(nil)
	branches := node.AvailableBranches()

	if branches != nil {
		t.Errorf("got: %+v, want: nil", branches)
	}
}
