package hoff

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/computestate"

	"github.com/rlespinasse/hoff/internal/utils"
	"github.com/rlespinasse/hoff/node"
)

func Test_NewComputation(t *testing.T) {
	var activatedSystem = NewNodeSystem()
	activatedSystem.Activate()

	var emptyContext = node.NewContextWithoutData()

	testCases := []struct {
		name                string
		givenSystem         *NodeSystem
		givenContext        *node.Context
		expectedComputation *Computation
		expectedError       error
	}{
		{
			name:                "Can't have a computation without system",
			givenSystem:         nil,
			givenContext:        emptyContext,
			expectedComputation: nil,
			expectedError:       errors.New("must have a node system to work properly"),
		},
		{
			name:                "Can't have a computation without activated system",
			givenSystem:         NewNodeSystem(),
			givenContext:        emptyContext,
			expectedComputation: nil,
			expectedError:       errors.New("must have an activated node system to work properly"),
		},
		{
			name:                "Can't have a computation without context",
			givenSystem:         activatedSystem,
			givenContext:        nil,
			expectedComputation: nil,
			expectedError:       errors.New("must have a context to work properly"),
		},
		{
			name:                "Can have a computation with validated node system and context",
			givenSystem:         activatedSystem,
			givenContext:        emptyContext,
			expectedComputation: &Computation{System: activatedSystem, Context: emptyContext},
			expectedError:       nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c, err := NewComputation(testCase.givenSystem, testCase.givenContext)

			if !cmp.Equal(c, testCase.expectedComputation) {
				t.Errorf("computation - got: %+v, want: %+v", c, testCase.expectedComputation)
			}
			if !cmp.Equal(err, testCase.expectedError, utils.ErrorComparator) {
				t.Errorf("error - got: %+v, want: %+v", err, testCase.expectedError)
			}
		})
	}
}

func Test_Computation_Compute(t *testing.T) {
	errorAction, _ := node.NewAction("errorAction", func(c *node.Context) error {
		return errors.New("action error")
	})
	errorDecision, _ := node.NewDecision("errorDecision", func(c *node.Context) (bool, error) {
		return false, errors.New("decision error")
	})
	writeAction, _ := node.NewAction("writeAction", func(c *node.Context) error {
		c.Store("write_action", "done")
		return nil
	})
	writeAnotherAction, _ := node.NewAction("writeAnotherAction", func(c *node.Context) error {
		c.Store("write_another_action", "done")
		return nil
	})
	readAction, _ := node.NewAction("readAction", func(c *node.Context) error {
		v, _ := c.Read("write_action")
		c.Store("read_action", fmt.Sprintf("the content of write_action is %v", v))
		return nil
	})
	deleteAnotherAction, _ := node.NewAction("deleteAnotherAction", func(c *node.Context) error {
		c.Delete("write_another_action")
		return nil
	})
	writeActionKeyIsPresent, _ := node.NewDecision("writeActionKeyIsPresent", func(c *node.Context) (bool, error) {
		return c.HaveKey("write_action"), nil
	})
	writeAnotherActionKeyIsPresent, _ := node.NewDecision("writeAnotherActionKeyIsPresent", func(c *node.Context) (bool, error) {
		return c.HaveKey("write_another_action"), nil
	})

	testCases := []struct {
		name                string
		givenNodes          []node.Node
		givenNodesJoinModes map[node.Node]JoinMode
		givenLinks          []nodeLink
		givenContextData    map[string]interface{}
		expectedStatus      bool
		expectedContextData map[string]interface{}
		expectedReport      map[node.Node]computestate.ComputeState
	}{
		{
			name:                "Can compute empty validated system",
			expectedStatus:      true,
			expectedContextData: map[string]interface{}{},
			expectedReport:      map[node.Node]computestate.ComputeState{},
		},
		{
			name: "Can compute one action node system",
			givenNodes: []node.Node{
				writeAction,
			},
			expectedStatus: true,
			expectedContextData: map[string]interface{}{
				"write_action": "done",
			},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAction: computestate.Continue(),
			},
		},
		{
			name: "Can compute 2 action nodes system",
			givenNodes: []node.Node{
				writeAction,
				writeAnotherAction,
			},
			expectedStatus: true,
			expectedContextData: map[string]interface{}{
				"write_action":         "done",
				"write_another_action": "done",
			},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAction:        computestate.Continue(),
				writeAnotherAction: computestate.Continue(),
			},
		},
		{
			name: "Can compute 2 linked action system (ordered declaration)",
			givenNodes: []node.Node{
				writeAction,
				readAction,
			},
			givenLinks: []nodeLink{
				newNodeLink(writeAction, readAction),
			},
			expectedStatus: true,
			expectedContextData: map[string]interface{}{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAction: computestate.Continue(),
				readAction:  computestate.Continue(),
			},
		},
		{
			name: "Can compute 2 linked action system (unordered declaration)",
			givenNodes: []node.Node{
				readAction,
				writeAction,
			},
			givenLinks: []nodeLink{
				newNodeLink(writeAction, readAction),
			},
			expectedStatus: true,
			expectedContextData: map[string]interface{}{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAction: computestate.Continue(),
				readAction:  computestate.Continue(),
			},
		},
		{
			name: "Can compute decision-based (branch 'true') system (ordered declaration)",
			givenNodes: []node.Node{
				writeAction,
				writeActionKeyIsPresent,
				readAction,
				deleteAnotherAction,
			},
			givenLinks: []nodeLink{
				newNodeLink(writeAction, writeActionKeyIsPresent),
				newNodeLinkOnBranch(writeActionKeyIsPresent, readAction, true),
				newNodeLinkOnBranch(writeActionKeyIsPresent, deleteAnotherAction, false),
			},
			expectedStatus: true,
			expectedContextData: map[string]interface{}{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAction:             computestate.Continue(),
				writeActionKeyIsPresent: computestate.ContinueOnBranch(true),
				readAction:              computestate.Continue(),
				deleteAnotherAction:     computestate.Skip(),
			},
		},
		{
			name: "Can compute decision-based (branch 'false') system (ordered declaration)",
			givenNodes: []node.Node{
				writeAnotherAction,
				writeActionKeyIsPresent,
				readAction,
				deleteAnotherAction,
			},
			givenLinks: []nodeLink{
				newNodeLink(writeAnotherAction, writeActionKeyIsPresent),
				newNodeLinkOnBranch(writeActionKeyIsPresent, readAction, true),
				newNodeLinkOnBranch(writeActionKeyIsPresent, deleteAnotherAction, false),
			},
			expectedStatus:      true,
			expectedContextData: map[string]interface{}{},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAnotherAction:      computestate.Continue(),
				writeActionKeyIsPresent: computestate.ContinueOnBranch(false),
				readAction:              computestate.Skip(),
				deleteAnotherAction:     computestate.Continue(),
			},
		},
		{
			name: "Can compute decision-based (branch 'true') system (unordered declaration)",
			givenNodes: []node.Node{
				deleteAnotherAction,
				readAction,
				writeActionKeyIsPresent,
				writeAction,
			},
			givenLinks: []nodeLink{
				newNodeLink(writeAction, writeActionKeyIsPresent),
				newNodeLinkOnBranch(writeActionKeyIsPresent, readAction, true),
				newNodeLinkOnBranch(writeActionKeyIsPresent, deleteAnotherAction, false),
			},
			expectedStatus: true,
			expectedContextData: map[string]interface{}{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAction:             computestate.Continue(),
				writeActionKeyIsPresent: computestate.ContinueOnBranch(true),
				readAction:              computestate.Continue(),
				deleteAnotherAction:     computestate.Skip(),
			},
		},
		{
			name: "Can compute decision-based (branch 'false') system (unordered declaration)",
			givenNodes: []node.Node{
				deleteAnotherAction,
				readAction,
				writeActionKeyIsPresent,
				writeAnotherAction,
			},
			givenLinks: []nodeLink{
				newNodeLink(writeAnotherAction, writeActionKeyIsPresent),
				newNodeLinkOnBranch(writeActionKeyIsPresent, readAction, true),
				newNodeLinkOnBranch(writeActionKeyIsPresent, deleteAnotherAction, false),
			},
			expectedStatus:      true,
			expectedContextData: map[string]interface{}{},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAnotherAction:      computestate.Continue(),
				writeActionKeyIsPresent: computestate.ContinueOnBranch(false),
				readAction:              computestate.Skip(),
				deleteAnotherAction:     computestate.Continue(),
			},
		},
		{
			name: "Can compute a node system with one erroring action node",
			givenNodes: []node.Node{
				writeAction,
				writeActionKeyIsPresent,
				errorAction,
				readAction,
				deleteAnotherAction,
			},
			givenLinks: []nodeLink{
				newNodeLink(writeAction, writeActionKeyIsPresent),
				newNodeLinkOnBranch(writeActionKeyIsPresent, errorAction, true),
				newNodeLink(errorAction, readAction),
				newNodeLinkOnBranch(writeActionKeyIsPresent, deleteAnotherAction, false),
			},
			expectedStatus: false,
			expectedContextData: map[string]interface{}{
				"write_action": "done",
			},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAction:             computestate.Continue(),
				writeActionKeyIsPresent: computestate.ContinueOnBranch(true),
				errorAction:             computestate.Abort(errors.New("action error")),
			},
		},
		{
			name: "Can compute another node system with one erroring action node",
			givenNodes: []node.Node{
				writeAnotherAction,
				writeActionKeyIsPresent,
				errorAction,
				readAction,
			},
			givenLinks: []nodeLink{
				newNodeLink(writeAnotherAction, writeActionKeyIsPresent),
				newNodeLinkOnBranch(writeActionKeyIsPresent, errorAction, false),
				newNodeLink(errorAction, readAction),
			},
			expectedStatus: false,
			expectedContextData: map[string]interface{}{
				"write_another_action": "done",
			},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAnotherAction:      computestate.Continue(),
				writeActionKeyIsPresent: computestate.ContinueOnBranch(false),
				errorAction:             computestate.Abort(errors.New("action error")),
			},
		},
		{
			name: "Can compute a node system with one erroring decision node",
			givenNodes: []node.Node{
				writeAction,
				errorDecision,
				readAction,
				deleteAnotherAction,
			},
			givenLinks: []nodeLink{
				newNodeLink(writeAction, errorDecision),
				newNodeLinkOnBranch(errorDecision, readAction, true),
				newNodeLinkOnBranch(errorDecision, deleteAnotherAction, false),
			},
			expectedStatus: false,
			expectedContextData: map[string]interface{}{
				"write_action": "done",
			},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAction:   computestate.Continue(),
				errorDecision: computestate.Abort(errors.New("decision error")),
			},
		},
		{
			name: "Can compute a node system with fork links",
			givenNodes: []node.Node{
				writeAction,
				readAction,
				errorAction,
				deleteAnotherAction,
			},
			givenLinks: []nodeLink{
				newNodeLink(writeAction, deleteAnotherAction),
				newNodeLink(writeAction, errorAction),
				newNodeLink(writeAction, readAction),
			},
			expectedStatus: false,
			expectedContextData: map[string]interface{}{
				"write_action": "done",
			},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAction:         computestate.Continue(),
				deleteAnotherAction: computestate.Continue(),
				errorAction:         computestate.Abort(errors.New("action error")),
			},
		},
		{
			name: "Can compute a node system with join and links",
			givenNodes: []node.Node{
				writeAction,
				writeAnotherAction,
				readAction,
				deleteAnotherAction,
			},
			givenNodesJoinModes: map[node.Node]JoinMode{
				readAction: JoinAnd,
			},
			givenLinks: []nodeLink{
				newNodeLink(writeAction, readAction),
				newNodeLink(writeAnotherAction, readAction),
				newNodeLink(readAction, deleteAnotherAction),
			},
			expectedStatus: true,
			expectedContextData: map[string]interface{}{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAction:         computestate.Continue(),
				writeAnotherAction:  computestate.Continue(),
				readAction:          computestate.Continue(),
				deleteAnotherAction: computestate.Continue(),
			},
		},
		{
			name: "Can compute a node system with partial join and links",
			givenNodes: []node.Node{
				writeAction,
				writeActionKeyIsPresent,
				writeAnotherAction,
				readAction,
				deleteAnotherAction,
			},
			givenNodesJoinModes: map[node.Node]JoinMode{
				readAction: JoinAnd,
			},
			givenLinks: []nodeLink{
				newNodeLink(writeAction, readAction),
				newNodeLinkOnBranch(writeActionKeyIsPresent, readAction, false),
				newNodeLink(writeAnotherAction, readAction),
				newNodeLink(readAction, deleteAnotherAction),
			},
			expectedStatus: true,
			expectedContextData: map[string]interface{}{
				"write_action":         "done",
				"write_another_action": "done",
			},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAction:             computestate.Continue(),
				writeActionKeyIsPresent: computestate.ContinueOnBranch(true),
				writeAnotherAction:      computestate.Continue(),
				readAction:              computestate.Skip(),
				deleteAnotherAction:     computestate.Skip(),
			},
		},
		{
			name: "Can compute a node system with join or links",
			givenNodes: []node.Node{
				writeAction,
				writeAnotherActionKeyIsPresent,
				readAction,
				deleteAnotherAction,
			},
			givenNodesJoinModes: map[node.Node]JoinMode{
				readAction: JoinOr,
			},
			givenLinks: []nodeLink{
				newNodeLink(writeAction, readAction),
				newNodeLinkOnBranch(writeAnotherActionKeyIsPresent, readAction, true),
				newNodeLink(readAction, deleteAnotherAction),
			},
			expectedStatus: true,
			expectedContextData: map[string]interface{}{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAction:                    computestate.Continue(),
				writeAnotherActionKeyIsPresent: computestate.ContinueOnBranch(false),
				readAction:                     computestate.Continue(),
				deleteAnotherAction:            computestate.Continue(),
			},
		},
		{
			name: "Can compute a node system with partial join or links",
			givenNodes: []node.Node{
				writeActionKeyIsPresent,
				writeAnotherActionKeyIsPresent,
				readAction,
				deleteAnotherAction,
			},
			givenNodesJoinModes: map[node.Node]JoinMode{
				readAction: JoinOr,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(writeActionKeyIsPresent, readAction, true),
				newNodeLinkOnBranch(writeAnotherActionKeyIsPresent, readAction, true),
				newNodeLink(readAction, deleteAnotherAction),
			},
			expectedStatus:      true,
			expectedContextData: map[string]interface{}{},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeActionKeyIsPresent:        computestate.ContinueOnBranch(false),
				writeAnotherActionKeyIsPresent: computestate.ContinueOnBranch(false),
				readAction:                     computestate.Skip(),
				deleteAnotherAction:            computestate.Skip(),
			},
		},
		{
			name: "Can compute a node system with join or links who generate error",
			givenNodes: []node.Node{
				writeAction,
				errorAction,
				readAction,
				deleteAnotherAction,
			},
			givenNodesJoinModes: map[node.Node]JoinMode{
				readAction: JoinOr,
			},
			givenLinks: []nodeLink{
				newNodeLink(writeAction, readAction),
				newNodeLink(errorAction, readAction),
				newNodeLink(readAction, deleteAnotherAction),
			},
			expectedStatus: false,
			expectedContextData: map[string]interface{}{
				"write_action": "done",
			},
			expectedReport: map[node.Node]computestate.ComputeState{
				writeAction: computestate.Continue(),
				errorAction: computestate.Abort(errors.New("action error")),
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			system := NewNodeSystem()
			loadNodeSystem(system, testCase.givenNodes, testCase.givenNodesJoinModes, testCase.givenLinks)

			_, errs := system.IsValid()
			if errs != nil {
				t.Errorf("validation errors - %+v\n", errs)
			}
			system.Activate()

			if testCase.givenContextData == nil {
				testCase.givenContextData = make(map[string]interface{})
			}
			context := node.NewContext(testCase.givenContextData)

			c, err := NewComputation(system, context)
			if err != nil {
				t.Errorf("can't compute: %+v", err)
				t.FailNow()
			}

			c.Compute()
			if !cmp.Equal(c.Status, testCase.expectedStatus) {
				t.Errorf("computation is done - got: %+v, want: %+v", c.Status, testCase.expectedStatus)
			}

			if testCase.expectedContextData == nil {
				testCase.expectedContextData = make(map[string]interface{})
			}
			expectedContext := node.NewContext(testCase.expectedContextData)
			if !cmp.Equal(c.Context, expectedContext) {
				t.Errorf("context data - got: %+v, want: %+v", c.Context, expectedContext)
			}
			if !cmp.Equal(c.Report, testCase.expectedReport, utils.ErrorComparator) {
				t.Errorf("report - got: %+v, want: %+v", c.Report, testCase.expectedReport)
			}
		})
	}
}

func Test_Github_Issue_11_JoinMode_AND(t *testing.T) {
	testGithubIssue11JoinMode(JoinAnd, t)
}

func Test_Github_Issue_11_JoinMode_OR(t *testing.T) {
	testGithubIssue11JoinMode(JoinOr, t)
}

func testGithubIssue11JoinMode(mode JoinMode, t *testing.T) {
	action1, _ := node.NewAction("action1", func(c *node.Context) error {
		data, _ := c.Read("run_order")
		newData := append(data.([]string), "action1")
		c.Store("run_order", newData)
		return nil
	})
	decision2, _ := node.NewDecision("decision2", func(c *node.Context) (bool, error) {
		data, _ := c.Read("run_order")
		newData := append(data.([]string), "decision2")
		c.Store("run_order", newData)
		return true, nil
	})
	action3, _ := node.NewAction("action3", func(c *node.Context) error {
		data, _ := c.Read("run_order")
		newData := append(data.([]string), "action3")
		c.Store("run_order", newData)
		return nil
	})
	action4, _ := node.NewAction("action4", func(c *node.Context) error {
		data, _ := c.Read("run_order")
		newData := append(data.([]string), "action4")
		c.Store("run_order", newData)
		return nil
	})
	action5, _ := node.NewAction("action5", func(c *node.Context) error {
		data, _ := c.Read("run_order")
		newData := append(data.([]string), "action5")
		c.Store("run_order", newData)
		return nil
	})

	data := map[string]interface{}{
		"run_order": []string{},
	}

	ns := NewNodeSystem()
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

	cp, _ := NewComputation(ns, node.NewContext(data))
	cp.Compute()

	resultData, _ := cp.Context.Read("run_order")
	expectedData := []string{"action1", "decision2", "action3", "action5"}

	if !cmp.Equal(resultData, expectedData) {
		t.Errorf("run order - got: %+v, want: %+v", resultData, expectedData)
	}
}
