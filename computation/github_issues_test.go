package computation

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/node"
	"github.com/rlespinasse/hoff/system"
	"github.com/rlespinasse/hoff/system/joinmode"
)

var (
	action1, _ = node.NewAction("action1", func(c *node.Context) error {
		data, _ := c.Read("run_order")
		newData := append(data.([]string), "action1")
		c.Store("run_order", newData)
		return nil
	})
	decision2, _ = node.NewDecision("decision2", func(c *node.Context) (bool, error) {
		data, _ := c.Read("run_order")
		newData := append(data.([]string), "decision2")
		c.Store("run_order", newData)
		return true, nil
	})
	action3, _ = node.NewAction("action3", func(c *node.Context) error {
		data, _ := c.Read("run_order")
		newData := append(data.([]string), "action3")
		c.Store("run_order", newData)
		return nil
	})
	action4, _ = node.NewAction("action4", func(c *node.Context) error {
		data, _ := c.Read("run_order")
		newData := append(data.([]string), "action4")
		c.Store("run_order", newData)
		return nil
	})
	action5, _ = node.NewAction("action5", func(c *node.Context) error {
		data, _ := c.Read("run_order")
		newData := append(data.([]string), "action5")
		c.Store("run_order", newData)
		return nil
	})
)

func Test_Github_Issue_11_JoinMode_AND(t *testing.T) {
	testRunOrderWithJoinMode(joinmode.AND, t)
}

func Test_Github_Issue_11_JoinMode_OR(t *testing.T) {
	testRunOrderWithJoinMode(joinmode.OR, t)
}

func testRunOrderWithJoinMode(mode joinmode.JoinMode, t *testing.T) {
	data := map[string]interface{}{
		"run_order": []string{},
	}

	ns := system.New()
	ns.AddNode(action1)
	ns.AddNode(decision2)
	ns.AddNode(action3)
	ns.AddNode(action4)
	ns.AddNode(action5)
	ns.AddLink(action1, decision2)
	ns.AddLinkOnBranch(decision2, action3, true)
	ns.AddLinkOnBranch(decision2, action4, false)
	ns.AddLink(action3, action5)
	ns.AddLink(action1, action5)
	ns.ConfigureJoinModeOnNode(action5, mode)
	ns.Activate()

	cp, _ := New(ns, node.NewContext(data))
	cp.Compute()

	resultData, _ := cp.Context.Read("run_order")
	expectedData := []string{"action1", "decision2", "action3", "action5"}

	if !cmp.Equal(resultData, expectedData) {
		t.Errorf("run order - got: %+v, want: %+v", resultData, expectedData)
	}
}
