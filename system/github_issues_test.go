package system

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/internal/utils"
	"github.com/rlespinasse/hoff/node"
	"github.com/rlespinasse/hoff/system/joinmode"
)

func Test_Github_Issue_10(t *testing.T) {
	action1, _ := node.NewAction("action1", func(c *node.Context) error {
		return nil
	})
	decision2, _ := node.NewDecision("decision2", func(c *node.Context) (bool, error) {
		return true, nil
	})
	decision3, _ := node.NewDecision("decision3", func(c *node.Context) (bool, error) {
		return true, nil
	})
	action4, _ := node.NewAction("action4", func(c *node.Context) error {
		return nil
	})

	ns := New()
	ns.AddNode(action1)
	ns.AddNode(decision2)
	ns.AddNode(decision3)
	ns.AddNode(action4)
	ns.AddLink(action1, decision2)
	ns.AddLink(action1, decision3)
	ns.AddLinkOnBranch(decision2, action4, false)
	ns.AddLinkOnBranch(decision3, action4, true)
	ns.ConfigureJoinModeOnNode(action4, joinmode.NONE)
	_, errs := ns.IsValid()

	expectedErrors := []error{
		errors.New("can't have multiple links (2) to the same node: action4 without join mode"),
	}

	if !cmp.Equal(errs, expectedErrors, utils.ErrorComparator) {
		t.Errorf("errors - got: %+v, want: %+v", errs, expectedErrors)
	}
}
