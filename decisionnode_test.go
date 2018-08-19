package flows

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_NewDecisionNode(t *testing.T) {
	testCases := []struct {
		name          string
		givenFunc     func(*Context) (bool, error)
		expectedError error
	}{
		{
			name:          "Can't create an decision node without function",
			expectedError: errors.New("can't create decision node without function"),
		},
		{
			name:          "Can create an decision node",
			givenFunc:     func(*Context) (bool, error) { return true, nil },
			expectedError: nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			node, err := NewDecisionNode(testCase.givenFunc)

			if !cmp.Equal(err, testCase.expectedError, equalOptionForError) {
				t.Errorf("error - got: %+v, want: %+v", err, testCase.expectedError)
			}
			if testCase.givenFunc == nil && node != nil {
				t.Errorf("action node - got: %+v, want: <nil>", node)
			}
		})
	}
}

func Test_DecisionNode_Compute(t *testing.T) {
	passingBranchTrueNode, _ := NewDecisionNode(func(*Context) (bool, error) { return true, nil })
	passingBranchFalseNode, _ := NewDecisionNode(func(*Context) (bool, error) { return false, nil })
	failingNode, _ := NewDecisionNode(func(*Context) (bool, error) { return false, errors.New("error") })
	tc := []NodeTestCase{
		{
			name:                 "Should Pass on Branch 'true'",
			givenNode:            passingBranchTrueNode,
			expectedComputeState: ComputeStateBranchPass("true"),
		},
		{
			name:                 "Should Pass on Branch 'false'",
			givenNode:            passingBranchFalseNode,
			expectedComputeState: ComputeStateBranchPass("false"),
		},
		{
			name:                 "Should Fail",
			givenNode:            failingNode,
			expectedComputeState: ComputeStateStopOnError(errors.New("error")),
		},
	}
	RunTestOnNode(t, tc)
}

func Test_DecisionNode_AvailableBranches(t *testing.T) {
	node, _ := NewDecisionNode(nil)
	branches := node.AvailableBranches()
	expectedBranches := []string{"true", "false"}

	if !cmp.Equal(branches, expectedBranches) {
		t.Errorf("got: %+v, want: %+v", branches, expectedBranches)
	}
}
