package flow

import (
	"errors"
	"testing"
)

func Test_ActionNode_Run(t *testing.T) {
	tc := []NodeTestCase{
		{
			"Pass",
			contextData{},
			NewActionNode(func(*Context) (bool, error) {
				return true, nil
			}),
			RunStatePass(),
			contextData{},
		},
		{
			"Stop",
			contextData{},
			NewActionNode(func(*Context) (bool, error) {
				return false, nil
			}),
			RunStateStop(),
			contextData{},
		},
		{
			"Fail",
			contextData{},
			NewActionNode(func(*Context) (bool, error) {
				return false, errors.New("error")
			}),
			RunStateFail(errors.New("error")),
			contextData{},
		},
	}
	RunTestOnNode(t, tc)
}

func Test_ActionNode_AvailableBranches(t *testing.T) {
	node := NewActionNode(nil)
	branches := node.AvailableBranches()

	if branches != nil {
		t.Errorf("got: %+v, want: nil", branches)
	}
}
