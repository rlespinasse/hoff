package flow

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_DecisionNode_Compute(t *testing.T) {
	tc := []NodeTestCase{
		{
			name:                 "Pass True",
			givenContextData:     contextData{"key": "value"},
			givenNode:            NewDecisionNode(func(*Context) (bool, error) { return true, nil }),
			expectedComputeState: ComputeStateBranchPass("true"),
			expectedContextData:  contextData{"key": "value"},
		},
		{
			name:                 "Pass False",
			givenContextData:     contextData{},
			givenNode:            NewDecisionNode(func(*Context) (bool, error) { return false, nil }),
			expectedComputeState: ComputeStateBranchPass("false"),
			expectedContextData:  contextData{},
		},
		{
			name:                 "Fail",
			givenContextData:     contextData{},
			givenNode:            NewDecisionNode(func(*Context) (bool, error) { return false, errors.New("error") }),
			expectedComputeState: ComputeStateFail(errors.New("error")),
			expectedContextData:  contextData{},
		},
	}
	ComputeTestOnNode(t, tc)
}

func Test_DecisionNode_AvailableBranches(t *testing.T) {
	node := NewDecisionNode(nil)
	branches := node.AvailableBranches()
	expectedBranches := []string{"true", "false"}

	if !cmp.Equal(branches, expectedBranches) {
		t.Errorf("got: %+v, want: %+v", branches, expectedBranches)
	}
}

func Test_isDecisionNode_true(t *testing.T) {
	node := NewDecisionNode(nil)

	if !isDecisionNode(node) {
		t.Error("got: false, want: true")
	}
}

func Test_isDecisionNode_false(t *testing.T) {
	node := &SomeNode{}

	if isDecisionNode(node) {
		t.Error("got: true, want: false")
	}
}
