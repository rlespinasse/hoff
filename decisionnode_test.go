package flowengine

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_DecisionNode_Run(t *testing.T) {
	tc := []FlowNodeTestCase{
		{
			"Pass True",
			contextData{
				"key": "value",
			},
			NewDecisionNode(func(*FlowContext) (bool, error) {
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
			NewDecisionNode(func(*FlowContext) (bool, error) {
				return false, nil
			}),
			RunStateBranchPass("false"),
			contextData{},
		},
		{
			"Fail",
			contextData{},
			NewDecisionNode(func(*FlowContext) (bool, error) {
				return false, errors.New("error")
			}),
			RunStateFail(errors.New("error")),
			contextData{},
		},
	}
	RunTestOnFlowNode(t, tc)
}

func Test_DecisionNode_AvailableBranches(t *testing.T) {
	node := NewDecisionNode(nil)
	branches := node.AvailableBranches()
	expectedBranches := AvailablesBranches("true", "false")

	if !cmp.Equal(branches, expectedBranches) {
		t.Errorf("got: %#v, want: %#v", branches, expectedBranches)
	}
}
