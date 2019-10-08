package hoff

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/internal/utils"
)

var (
	continueNode, _ = NewActionNode("continueNode", func(*Context) error { return nil })
	abortNode, _    = NewActionNode("abortNode", func(*Context) error { return errors.New("error") })
)

func Test_NewActionNode(t *testing.T) {
	testCases := []struct {
		name          string
		givenFunc     func(*Context) error
		expectedError error
	}{
		{
			name:          "Can't create an action node without function",
			expectedError: errors.New("can't create action node without function"),
		},
		{
			name:          "Can create an action node",
			givenFunc:     func(*Context) error { return nil },
			expectedError: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			node, err := NewActionNode("ActionNode", testCase.givenFunc)

			if !cmp.Equal(err, testCase.expectedError, utils.ErrorComparator) {
				t.Errorf("error - got: %+v, want: %+v", err, testCase.expectedError)
			}
			if testCase.givenFunc == nil && node != nil {
				t.Errorf("action node - got: %+v, want: <nil>", node)
			}
		})
	}
}

func Test_ActionNode_Compute(t *testing.T) {
	tc := []NodeTestCase{
		{
			name:                 "Should Continue",
			givenNode:            continueNode,
			expectedComputeState: NewContinueComputeState(),
		},
		{
			name:                 "Should Abort",
			givenNode:            abortNode,
			expectedComputeState: NewAbortComputeState(errors.New("error")),
		},
	}
	RunTestOnNode(t, tc)
}

func Test_ActionNode_DecideCapability(t *testing.T) {
	if continueNode.DecideCapability() {
		t.Error("action node must have no decide capability")
	}
}

func Test_ActionNode_String(t *testing.T) {
	if continueNode.String() != "continueNode" {
		t.Error("action node must print its name")
	}
}
