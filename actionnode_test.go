package flow

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_NewActionNode(t *testing.T) {
	testCases := []struct {
		name          string
		givenFunc     func(*Context) (bool, error)
		expectedError error
	}{
		{
			name:          "Can't create an action node without function",
			expectedError: errors.New("can't create action node without function"),
		},
		{
			name:          "Can create an action node",
			givenFunc:     func(*Context) (bool, error) { return true, nil },
			expectedError: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			node, err := NewActionNode(testCase.givenFunc)

			if !cmp.Equal(err, testCase.expectedError, equalOptionForError) {
				t.Errorf("error - got: %+v, want: %+v", err, testCase.expectedError)
			}
			if testCase.givenFunc == nil && node != nil {
				t.Errorf("action node - got: %+v, want: <nil>", node)
			}
		})
	}
}

func Test_ActionNode_Compute(t *testing.T) {
	passingNode, _ := NewActionNode(func(*Context) (bool, error) { return true, nil })
	stoppingNode, _ := NewActionNode(func(*Context) (bool, error) { return false, nil })
	failingNode, _ := NewActionNode(func(*Context) (bool, error) { return false, errors.New("error") })
	tc := []NodeTestCase{
		{
			name:                 "Should Pass",
			givenNode:            passingNode,
			expectedComputeState: ComputeStatePass(),
		},
		{
			name:                 "Should Stop",
			givenNode:            stoppingNode,
			expectedComputeState: ComputeStateStop(),
		},
		{
			name:                 "Should Fail",
			givenNode:            failingNode,
			expectedComputeState: ComputeStateFail(errors.New("error")),
		},
	}
	RunTestOnNode(t, tc)
}

func Test_ActionNode_AvailableBranches(t *testing.T) {
	node, _ := NewActionNode(nil)
	branches := node.AvailableBranches()

	if branches != nil {
		t.Errorf("got: %+v, want: nil", branches)
	}
}
