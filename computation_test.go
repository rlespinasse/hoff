package hoff

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/rlespinasse/hoff/internal/utils"
)

func Test_NewComputation(t *testing.T) {
	var activatedSystem = NewNodeSystem()
	activatedSystem.Activate()

	var emptyContext = NewContextWithoutData()

	testCases := []struct {
		name                string
		givenSystem         *NodeSystem
		givenContext        *Context
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
	errorAction, _ := NewActionNode("errorAction", func(c *Context) error {
		return errors.New("action error")
	})
	errorDecision, _ := NewDecisionNode("errorDecision", func(c *Context) (bool, error) {
		return false, errors.New("decision error")
	})
	writeAction, _ := NewActionNode("writeAction", func(c *Context) error {
		c.Store("write_action", "done")
		return nil
	})
	writeAnotherAction, _ := NewActionNode("writeAnotherAction", func(c *Context) error {
		c.Store("write_another_action", "done")
		return nil
	})
	readAction, _ := NewActionNode("readAction", func(c *Context) error {
		v, _ := c.Read("write_action")
		c.Store("read_action", fmt.Sprintf("the content of write_action is %v", v))
		return nil
	})
	deleteAnotherAction, _ := NewActionNode("deleteAnotherAction", func(c *Context) error {
		c.Delete("write_another_action")
		return nil
	})
	writeActionKeyIsPresent, _ := NewDecisionNode("writeActionKeyIsPresent", func(c *Context) (bool, error) {
		return c.HaveKey("write_action"), nil
	})
	writeAnotherActionKeyIsPresent, _ := NewDecisionNode("writeAnotherActionKeyIsPresent", func(c *Context) (bool, error) {
		return c.HaveKey("write_another_action"), nil
	})

	testCases := []struct {
		name                string
		givenNodes          []Node
		givenNodesJoinModes map[Node]JoinMode
		givenLinks          []nodeLink
		givenContextData    map[string]interface{}
		expectedStatus      bool
		expectedContextData map[string]interface{}
		expectedReport      map[Node]ComputeState
	}{
		{
			name:                "Can compute empty validated system",
			expectedStatus:      true,
			expectedContextData: map[string]interface{}{},
			expectedReport:      map[Node]ComputeState{},
		},
		{
			name: "Can compute one action node system",
			givenNodes: []Node{
				writeAction,
			},
			expectedStatus: true,
			expectedContextData: map[string]interface{}{
				"write_action": "done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction: NewContinueComputeState(),
			},
		},
		{
			name: "Can compute 2 action nodes system",
			givenNodes: []Node{
				writeAction,
				writeAnotherAction,
			},
			expectedStatus: true,
			expectedContextData: map[string]interface{}{
				"write_action":         "done",
				"write_another_action": "done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction:        NewContinueComputeState(),
				writeAnotherAction: NewContinueComputeState(),
			},
		},
		{
			name: "Can compute 2 linked action system (ordered declaration)",
			givenNodes: []Node{
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
			expectedReport: map[Node]ComputeState{
				writeAction: NewContinueComputeState(),
				readAction:  NewContinueComputeState(),
			},
		},
		{
			name: "Can compute 2 linked action system (unordered declaration)",
			givenNodes: []Node{
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
			expectedReport: map[Node]ComputeState{
				writeAction: NewContinueComputeState(),
				readAction:  NewContinueComputeState(),
			},
		},
		{
			name: "Can compute decision-based (branch 'true') system (ordered declaration)",
			givenNodes: []Node{
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
			expectedReport: map[Node]ComputeState{
				writeAction:             NewContinueComputeState(),
				writeActionKeyIsPresent: NewContinueOnBranchComputeState(true),
				readAction:              NewContinueComputeState(),
				deleteAnotherAction:     NewSkipComputeState(),
			},
		},
		{
			name: "Can compute decision-based (branch 'false') system (ordered declaration)",
			givenNodes: []Node{
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
			expectedReport: map[Node]ComputeState{
				writeAnotherAction:      NewContinueComputeState(),
				writeActionKeyIsPresent: NewContinueOnBranchComputeState(false),
				readAction:              NewSkipComputeState(),
				deleteAnotherAction:     NewContinueComputeState(),
			},
		},
		{
			name: "Can compute decision-based (branch 'true') system (unordered declaration)",
			givenNodes: []Node{
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
			expectedReport: map[Node]ComputeState{
				writeAction:             NewContinueComputeState(),
				writeActionKeyIsPresent: NewContinueOnBranchComputeState(true),
				readAction:              NewContinueComputeState(),
				deleteAnotherAction:     NewSkipComputeState(),
			},
		},
		{
			name: "Can compute decision-based (branch 'false') system (unordered declaration)",
			givenNodes: []Node{
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
			expectedReport: map[Node]ComputeState{
				writeAnotherAction:      NewContinueComputeState(),
				writeActionKeyIsPresent: NewContinueOnBranchComputeState(false),
				readAction:              NewSkipComputeState(),
				deleteAnotherAction:     NewContinueComputeState(),
			},
		},
		{
			name: "Can compute a node system with one erroring action node",
			givenNodes: []Node{
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
			expectedReport: map[Node]ComputeState{
				writeAction:             NewContinueComputeState(),
				writeActionKeyIsPresent: NewContinueOnBranchComputeState(true),
				errorAction:             NewAbortComputeState(errors.New("action error")),
			},
		},
		{
			name: "Can compute another node system with one erroring action node",
			givenNodes: []Node{
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
			expectedReport: map[Node]ComputeState{
				writeAnotherAction:      NewContinueComputeState(),
				writeActionKeyIsPresent: NewContinueOnBranchComputeState(false),
				errorAction:             NewAbortComputeState(errors.New("action error")),
			},
		},
		{
			name: "Can compute a node system with one erroring decision node",
			givenNodes: []Node{
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
			expectedReport: map[Node]ComputeState{
				writeAction:   NewContinueComputeState(),
				errorDecision: NewAbortComputeState(errors.New("decision error")),
			},
		},
		{
			name: "Can compute a node system with fork links",
			givenNodes: []Node{
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
			expectedReport: map[Node]ComputeState{
				writeAction:         NewContinueComputeState(),
				deleteAnotherAction: NewContinueComputeState(),
				errorAction:         NewAbortComputeState(errors.New("action error")),
			},
		},
		{
			name: "Can compute a node system with join and links",
			givenNodes: []Node{
				writeAction,
				writeAnotherAction,
				readAction,
				deleteAnotherAction,
			},
			givenNodesJoinModes: map[Node]JoinMode{
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
			expectedReport: map[Node]ComputeState{
				writeAction:         NewContinueComputeState(),
				writeAnotherAction:  NewContinueComputeState(),
				readAction:          NewContinueComputeState(),
				deleteAnotherAction: NewContinueComputeState(),
			},
		},
		{
			name: "Can compute a node system with partial join and links",
			givenNodes: []Node{
				writeAction,
				writeActionKeyIsPresent,
				writeAnotherAction,
				readAction,
				deleteAnotherAction,
			},
			givenNodesJoinModes: map[Node]JoinMode{
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
			expectedReport: map[Node]ComputeState{
				writeAction:             NewContinueComputeState(),
				writeActionKeyIsPresent: NewContinueOnBranchComputeState(true),
				writeAnotherAction:      NewContinueComputeState(),
				readAction:              NewSkipComputeState(),
				deleteAnotherAction:     NewSkipComputeState(),
			},
		},
		{
			name: "Can compute a node system with join or links",
			givenNodes: []Node{
				writeAction,
				writeAnotherActionKeyIsPresent,
				readAction,
				deleteAnotherAction,
			},
			givenNodesJoinModes: map[Node]JoinMode{
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
			expectedReport: map[Node]ComputeState{
				writeAction:                    NewContinueComputeState(),
				writeAnotherActionKeyIsPresent: NewContinueOnBranchComputeState(false),
				readAction:                     NewContinueComputeState(),
				deleteAnotherAction:            NewContinueComputeState(),
			},
		},
		{
			name: "Can compute a node system with partial join or links",
			givenNodes: []Node{
				writeActionKeyIsPresent,
				writeAnotherActionKeyIsPresent,
				readAction,
				deleteAnotherAction,
			},
			givenNodesJoinModes: map[Node]JoinMode{
				readAction: JoinOr,
			},
			givenLinks: []nodeLink{
				newNodeLinkOnBranch(writeActionKeyIsPresent, readAction, true),
				newNodeLinkOnBranch(writeAnotherActionKeyIsPresent, readAction, true),
				newNodeLink(readAction, deleteAnotherAction),
			},
			expectedStatus:      true,
			expectedContextData: map[string]interface{}{},
			expectedReport: map[Node]ComputeState{
				writeActionKeyIsPresent:        NewContinueOnBranchComputeState(false),
				writeAnotherActionKeyIsPresent: NewContinueOnBranchComputeState(false),
				readAction:                     NewSkipComputeState(),
				deleteAnotherAction:            NewSkipComputeState(),
			},
		},
		{
			name: "Can compute a node system with join or links who generate error",
			givenNodes: []Node{
				writeAction,
				errorAction,
				readAction,
				deleteAnotherAction,
			},
			givenNodesJoinModes: map[Node]JoinMode{
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
			expectedReport: map[Node]ComputeState{
				writeAction: NewContinueComputeState(),
				errorAction: NewAbortComputeState(errors.New("action error")),
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
			context := NewContext(testCase.givenContextData)

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
			expectedContext := NewContext(testCase.expectedContextData)
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
	action1, _ := NewActionNode("action1", func(c *Context) error {
		data, _ := c.Read("run_order")
		newData := append(data.([]string), "action1")
		c.Store("run_order", newData)
		return nil
	})
	decision2, _ := NewDecisionNode("decision2", func(c *Context) (bool, error) {
		data, _ := c.Read("run_order")
		newData := append(data.([]string), "decision2")
		c.Store("run_order", newData)
		return true, nil
	})
	action3, _ := NewActionNode("action3", func(c *Context) error {
		data, _ := c.Read("run_order")
		newData := append(data.([]string), "action3")
		c.Store("run_order", newData)
		return nil
	})
	action4, _ := NewActionNode("action4", func(c *Context) error {
		data, _ := c.Read("run_order")
		newData := append(data.([]string), "action4")
		c.Store("run_order", newData)
		return nil
	})
	action5, _ := NewActionNode("action5", func(c *Context) error {
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

	cp, _ := NewComputation(ns, NewContext(data))
	cp.Compute()

	resultData, _ := cp.Context.Read("run_order")
	expectedData := []string{"action1", "decision2", "action3", "action5"}

	if !cmp.Equal(resultData, expectedData) {
		t.Errorf("run order - got: %+v, want: %+v", resultData, expectedData)
	}
}
