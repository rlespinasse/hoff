package flow

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_ActionNode_Run(t *testing.T) {
	tc := []FlowNodeTestCase{
		{
			"Pass",
			contextData{},
			NewActionNode(func(*FlowContext) (bool, error) {
				return true, nil
			}),
			RunStatePass(),
			contextData{},
		},
		{
			"Stop",
			contextData{},
			NewActionNode(func(*FlowContext) (bool, error) {
				return false, nil
			}),
			RunStateStop(),
			contextData{},
		},
		{
			"Fail",
			contextData{},
			NewActionNode(func(*FlowContext) (bool, error) {
				return false, errors.New("error")
			}),
			RunStateFail(errors.New("error")),
			contextData{},
		},
	}
	RunTestOnFlowNode(t, tc)
}

func Test_ActionNode_AvailableBranches(t *testing.T) {
	node := NewActionNode(nil)
	branches := node.AvailableBranches()
	expectedNodeBranches := AvailablesBranches()

	if !cmp.Equal(branches, expectedNodeBranches) {
		t.Errorf("got: %#v, want: %#v", branches, expectedNodeBranches)
	}
}
