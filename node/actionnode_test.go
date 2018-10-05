package node

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/computestate"
	"github.com/rlespinasse/hoff/internal/utils"
)

var (
	continueNode, _ = NewAction("continueNode", func(*Context) (bool, error) { return true, nil })
	stopNode, _     = NewAction("stopNode", func(*Context) (bool, error) { return false, nil })
	abortNode, _    = NewAction("abortNode", func(*Context) (bool, error) { return false, errors.New("error") })
)

func Test_NewAction(t *testing.T) {
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
			node, err := NewAction("ActionNode", testCase.givenFunc)

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
			expectedComputeState: computestate.Continue(),
		},
		{
			name:                 "Should Stop",
			givenNode:            stopNode,
			expectedComputeState: computestate.Stop(),
		},
		{
			name:                 "Should Abort",
			givenNode:            abortNode,
			expectedComputeState: computestate.Abort(errors.New("error")),
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
