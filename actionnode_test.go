package hoff

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
			node, err := NewActionNode("ActionNode", testCase.givenFunc)

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
	continueNode, _ := NewActionNode("continueNode", func(*Context) (bool, error) { return true, nil })
	stopNode, _ := NewActionNode("stopNode", func(*Context) (bool, error) { return false, nil })
	abortNode, _ := NewActionNode("abortNode", func(*Context) (bool, error) { return false, errors.New("error") })
	tc := []NodeTestCase{
		{
			name:                 "Should Continue",
			givenNode:            continueNode,
			expectedComputeState: ComputeStateContinue(),
		},
		{
			name:                 "Should Stop",
			givenNode:            stopNode,
			expectedComputeState: ComputeStateStop(),
		},
		{
			name:                 "Should Abort",
			givenNode:            abortNode,
			expectedComputeState: ComputeStateAbort(errors.New("error")),
		},
	}
	RunTestOnNode(t, tc)
}
