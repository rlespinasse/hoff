package system

import (
	"errors"
	"fmt"
	"testing"

	"github.com/rlespinasse/hoff/internal/nodelink"

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

func Test_Github_Issue_16(t *testing.T) {
	trigger, _ := node.NewDecision("trigger", func(c *node.Context) (bool, error) {
		return true, nil
	})
	a1, _ := node.NewAction("a1", func(c *node.Context) error {
		return nil
	})
	a2, _ := node.NewAction("a2", func(c *node.Context) error {
		return nil
	})
	a3, _ := node.NewAction("a3", func(c *node.Context) error {
		return nil
	})

	ns := New()
	ns.AddNode(trigger)
	ns.AddNode(a1)
	ns.AddNode(a2)
	ns.AddNode(a3)
	ns.AddLinkOnBranch(trigger, a1, true)
	ns.AddLink(a1, a2)
	ns.AddLink(a2, a3)
	ns.AddLink(a3, a2)
	ns.ConfigureJoinModeOnNode(a2, joinmode.AND)
	_, errs := ns.IsValid()

	expectedErrors := []error{
		fmt.Errorf("Can't have cycle in links between nodes: %+v", []nodelink.NodeLink{
			nodelink.New(a2, a3),
			nodelink.New(a3, a2),
		}),
	}

	if !cmp.Equal(errs, expectedErrors, utils.ErrorComparator) {
		t.Errorf("errors - got: %+v, want: %+v", errs, expectedErrors)
	}
}

func Test_multiple_cycles(t *testing.T) {
	trigger, _ := node.NewDecision("trigger", func(c *node.Context) (bool, error) {
		return true, nil
	})
	a1, _ := node.NewAction("a1", func(c *node.Context) error {
		return nil
	})
	a2, _ := node.NewAction("a2", func(c *node.Context) error {
		return nil
	})
	a3, _ := node.NewAction("a3", func(c *node.Context) error {
		return nil
	})
	a4, _ := node.NewAction("a4", func(c *node.Context) error {
		return nil
	})
	a5, _ := node.NewAction("a5", func(c *node.Context) error {
		return nil
	})
	a6, _ := node.NewAction("a6", func(c *node.Context) error {
		return nil
	})
	a7, _ := node.NewAction("a7", func(c *node.Context) error {
		return nil
	})

	ns := New()
	ns.AddNode(trigger)
	ns.AddNode(a1)
	ns.AddNode(a2)
	ns.AddNode(a3)
	ns.AddNode(a4)
	ns.AddNode(a5)
	ns.AddNode(a6)
	ns.AddNode(a7)

	// Cycle between a2 and a3
	ns.AddLinkOnBranch(trigger, a1, true)
	ns.AddLink(a1, a2)
	ns.AddLink(a2, a3)
	ns.AddLink(a3, a2)
	ns.ConfigureJoinModeOnNode(a2, joinmode.AND)

	// Cycle between a5, a6, and a7
	ns.AddLinkOnBranch(trigger, a4, false)
	ns.AddLink(a4, a5)
	ns.AddLink(a5, a6)
	ns.AddLink(a6, a7)
	ns.AddLink(a7, a5)
	ns.ConfigureJoinModeOnNode(a5, joinmode.AND)
	_, errs := ns.IsValid()

	expectedErrors := []error{
		fmt.Errorf("Can't have cycle in links between nodes: %+v", []nodelink.NodeLink{
			nodelink.New(a2, a3),
			nodelink.New(a3, a2),
		}),
		fmt.Errorf("Can't have cycle in links between nodes: %+v", []nodelink.NodeLink{
			nodelink.New(a5, a6),
			nodelink.New(a6, a7),
			nodelink.New(a7, a5),
		}),
	}

	if !cmp.Equal(errs, expectedErrors, utils.ErrorComparator) {
		t.Errorf("errors - got: %+v, want: %+v", errs, expectedErrors)
	}
}
