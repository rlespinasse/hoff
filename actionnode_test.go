package flow

import (
	"errors"
	"testing"
)

func Test_ActionNode_Compute(t *testing.T) {
	tc := []NodeTestCase{
		{
			"Pass",
			contextData{},
			NewActionNode(func(*Context) (bool, error) {
				return true, nil
			}),
			ComputeStatePass(),
			contextData{},
		},
		{
			"Stop",
			contextData{},
			NewActionNode(func(*Context) (bool, error) {
				return false, nil
			}),
			ComputeStateStop(),
			contextData{},
		},
		{
			"Fail",
			contextData{},
			NewActionNode(func(*Context) (bool, error) {
				return false, errors.New("error")
			}),
			ComputeStateFail(errors.New("error")),
			contextData{},
		},
	}
	ComputeTestOnNode(t, tc)
}

func Test_ActionNode_AvailableBranches(t *testing.T) {
	node := NewActionNode(nil)
	branches := node.AvailableBranches()

	if branches != nil {
		t.Errorf("got: %+v, want: nil", branches)
	}
}
