package namingishard

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
			node, err := NewDecisionNode("DecisionNode", testCase.givenFunc)

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
	passingBranchTrueNode, _ := NewDecisionNode("passingBranchTrueNode", func(*Context) (bool, error) { return true, nil })
	passingBranchFalseNode, _ := NewDecisionNode("passingBranchFalseNode", func(*Context) (bool, error) { return false, nil })
	failingNode, _ := NewDecisionNode("failingNode", func(*Context) (bool, error) { return false, errors.New("error") })
	tc := []NodeTestCase{
		{
			name:                 "Should Pass on Branch 'true'",
			givenNode:            passingBranchTrueNode,
			expectedComputeState: ComputeStateContinueOnBranch(true),
		},
		{
			name:                 "Should Pass on Branch 'false'",
			givenNode:            passingBranchFalseNode,
			expectedComputeState: ComputeStateContinueOnBranch(false),
		},
		{
			name:                 "Should Fail",
			givenNode:            failingNode,
			expectedComputeState: ComputeStateAbort(errors.New("error")),
		},
	}
	RunTestOnNode(t, tc)
}
