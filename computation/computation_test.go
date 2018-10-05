package computation

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/computestate"

	"github.com/rlespinasse/hoff/internal/nodelink"
	"github.com/rlespinasse/hoff/internal/utils"
	"github.com/rlespinasse/hoff/node"
	"github.com/rlespinasse/hoff/system"
	"github.com/rlespinasse/hoff/system/joinmode"
)

func Test_New(t *testing.T) {
	var activatedSystem = system.New()
	activatedSystem.Activate()

	var emptyContext = node.NewContextWithoutData()

	testCases := []struct {
		name                string
		givenSystem         *system.NodeSystem
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
			givenSystem:         system.New(),
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
			c, err := New(testCase.givenSystem, testCase.givenContext)

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
	errorAction, _ := node.NewAction("errorAction", func(c *node.Context) (bool, error) {
		return true, errors.New("action error")
	})
	errorDecision, _ := node.NewDecision("errorDecision", func(c *node.Context) (bool, error) {
		return false, errors.New("decision error")
	})
	writeAction, _ := node.NewAction("writeAction", func(c *node.Context) (bool, error) {
		c.Store("write_action", "done")
		return true, nil
	})
	writeAnotherAction, _ := node.NewAction("writeAnotherAction", func(c *node.Context) (bool, error) {
		c.Store("write_another_action", "done")
		return true, nil
	})
	readAction, _ := node.NewAction("readAction", func(c *node.Context) (bool, error) {
		v, ok := c.Read("write_action")
		if !ok {
			return false, nil
		}
		c.Store("read_action", fmt.Sprintf("the content of write_action is %v", v))
		return true, nil
	})
	deleteAnotherAction, _ := node.NewAction("deleteAnotherAction", func(c *node.Context) (bool, error) {
		c.Delete("write_another_action")
		return true, nil
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
		givenNodesJoinModes map[node.Node]joinmode.JoinMode
		givenLinks          []nodelink.NodeLink
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
			name: "Can compute one stop action node system",
			givenNodes: []node.Node{
				readAction,
			},
			expectedStatus:      true,
			expectedContextData: map[string]interface{}{},
			expectedReport: map[node.Node]computestate.ComputeState{
				readAction: computestate.Stop(),
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
			givenLinks: []nodelink.NodeLink{
				nodelink.New(writeAction, readAction),
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
			givenLinks: []nodelink.NodeLink{
				nodelink.New(writeAction, readAction),
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
			givenLinks: []nodelink.NodeLink{
				nodelink.New(writeAction, writeActionKeyIsPresent),
				nodelink.NewOnBranch(writeActionKeyIsPresent, readAction, true),
				nodelink.NewOnBranch(writeActionKeyIsPresent, deleteAnotherAction, false),
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
			givenLinks: []nodelink.NodeLink{
				nodelink.New(writeAnotherAction, writeActionKeyIsPresent),
				nodelink.NewOnBranch(writeActionKeyIsPresent, readAction, true),
				nodelink.NewOnBranch(writeActionKeyIsPresent, deleteAnotherAction, false),
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
			givenLinks: []nodelink.NodeLink{
				nodelink.New(writeAction, writeActionKeyIsPresent),
				nodelink.NewOnBranch(writeActionKeyIsPresent, readAction, true),
				nodelink.NewOnBranch(writeActionKeyIsPresent, deleteAnotherAction, false),
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
			givenLinks: []nodelink.NodeLink{
				nodelink.New(writeAnotherAction, writeActionKeyIsPresent),
				nodelink.NewOnBranch(writeActionKeyIsPresent, readAction, true),
				nodelink.NewOnBranch(writeActionKeyIsPresent, deleteAnotherAction, false),
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
			givenLinks: []nodelink.NodeLink{
				nodelink.New(writeAction, writeActionKeyIsPresent),
				nodelink.NewOnBranch(writeActionKeyIsPresent, errorAction, true),
				nodelink.New(errorAction, readAction),
				nodelink.NewOnBranch(writeActionKeyIsPresent, deleteAnotherAction, false),
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
			givenLinks: []nodelink.NodeLink{
				nodelink.New(writeAnotherAction, writeActionKeyIsPresent),
				nodelink.NewOnBranch(writeActionKeyIsPresent, errorAction, false),
				nodelink.New(errorAction, readAction),
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
			givenLinks: []nodelink.NodeLink{
				nodelink.New(writeAction, errorDecision),
				nodelink.NewOnBranch(errorDecision, readAction, true),
				nodelink.NewOnBranch(errorDecision, deleteAnotherAction, false),
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
			givenLinks: []nodelink.NodeLink{
				nodelink.New(writeAction, deleteAnotherAction),
				nodelink.New(writeAction, errorAction),
				nodelink.New(writeAction, readAction),
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
			givenNodesJoinModes: map[node.Node]joinmode.JoinMode{
				readAction: joinmode.AND,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.New(writeAction, readAction),
				nodelink.New(writeAnotherAction, readAction),
				nodelink.New(readAction, deleteAnotherAction),
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
			givenNodesJoinModes: map[node.Node]joinmode.JoinMode{
				readAction: joinmode.AND,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.New(writeAction, readAction),
				nodelink.NewOnBranch(writeActionKeyIsPresent, readAction, false),
				nodelink.New(writeAnotherAction, readAction),
				nodelink.New(readAction, deleteAnotherAction),
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
			givenNodesJoinModes: map[node.Node]joinmode.JoinMode{
				readAction: joinmode.OR,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.New(writeAction, readAction),
				nodelink.NewOnBranch(writeAnotherActionKeyIsPresent, readAction, true),
				nodelink.New(readAction, deleteAnotherAction),
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
			givenNodesJoinModes: map[node.Node]joinmode.JoinMode{
				readAction: joinmode.OR,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.NewOnBranch(writeActionKeyIsPresent, readAction, true),
				nodelink.NewOnBranch(writeAnotherActionKeyIsPresent, readAction, true),
				nodelink.New(readAction, deleteAnotherAction),
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
			givenNodesJoinModes: map[node.Node]joinmode.JoinMode{
				readAction: joinmode.OR,
			},
			givenLinks: []nodelink.NodeLink{
				nodelink.New(writeAction, readAction),
				nodelink.New(errorAction, readAction),
				nodelink.New(readAction, deleteAnotherAction),
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
			system := system.New()
			load(system, testCase.givenNodes, testCase.givenNodesJoinModes, testCase.givenLinks)

			_, errs := system.IsValid()
			if errs != nil {
				t.Errorf("validation errors - %+v\n", errs)
			}
			system.Activate()

			if testCase.givenContextData == nil {
				testCase.givenContextData = make(map[string]interface{})
			}
			context := node.NewContext(testCase.givenContextData)

			c, err := New(system, context)
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

func load(system *system.NodeSystem, nodes []node.Node, nodesJoinModes map[node.Node]joinmode.JoinMode, links []nodelink.NodeLink) []error {
	var errs []error
	for _, node := range nodes {
		_, err := system.AddNode(node)
		if err != nil {
			errs = append(errs, err)
		}
	}
	for node, mode := range nodesJoinModes {
		_, err := system.ConfigureJoinModeOnNode(node, mode)
		if err != nil {
			errs = append(errs, err)
		}
	}
	for _, link := range links {
		if link.Branch == nil {
			_, err := system.AddLink(link.From, link.To)
			if err != nil {
				errs = append(errs, err)
			}
		} else {
			_, err := system.AddLinkOnBranch(link.From, link.To, *link.Branch)
			if err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errs
}
