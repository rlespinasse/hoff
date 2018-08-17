package flow

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_DecisionNode_Run(t *testing.T) {
	tc := []NodeTestCase{
		{
			"Pass True",
			contextData{
				"key": "value",
			},
			NewDecisionNode(func(*Context) (bool, error) {
				return true, nil
			}),
			RunStateBranchPass("true"),
			contextData{
				"key": "value",
			},
		},
		{
			"Pass False",
			contextData{},
			NewDecisionNode(func(*Context) (bool, error) {
				return false, nil
			}),
			RunStateBranchPass("false"),
			contextData{},
		},
		{
			"Fail",
			contextData{},
			NewDecisionNode(func(*Context) (bool, error) {
				return false, errors.New("error")
			}),
			RunStateFail(errors.New("error")),
			contextData{},
		},
	}
	RunTestOnNode(t, tc)
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
