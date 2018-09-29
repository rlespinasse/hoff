package hoff

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_NewComputation(t *testing.T) {
	var validatedSystem = NewNodeSystem()
	validatedSystem.validity = true

	var emptyContext = NewContext()

	testCases := []struct {
		name                string
		givenSystem         *NodeSystem
		givenContext        *Context
		expectedComputation *computation
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
			name:                "Can't have a computation without validated system",
			givenSystem:         NewNodeSystem(),
			givenContext:        emptyContext,
			expectedComputation: nil,
			expectedError:       errors.New("must have a validated node system to work properly"),
		},
		{
			name:                "Can't have a computation without context",
			givenSystem:         validatedSystem,
			givenContext:        nil,
			expectedComputation: nil,
			expectedError:       errors.New("must have a context to work properly"),
		},
		{
			name:                "Can have a computation with validated node system and context",
			givenSystem:         validatedSystem,
			givenContext:        emptyContext,
			expectedComputation: &computation{computation: false, system: validatedSystem, context: emptyContext},
			expectedError:       nil,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c, err := NewComputation(testCase.givenSystem, testCase.givenContext)

			if !cmp.Equal(c, testCase.expectedComputation) {
				t.Errorf("computation - got: %+v, want: %+v", c, testCase.expectedComputation)
			}
			if !cmp.Equal(err, testCase.expectedError, equalOptionForError) {
				t.Errorf("error - got: %+v, want: %+v", err, testCase.expectedError)
			}
		})
	}
}

func Test_Computation_Compute(t *testing.T) {
	errorAction, _ := NewActionNode("errorAction", func(c *Context) (bool, error) {
		return true, errors.New("action error")
	})
	errorDecision, _ := NewDecisionNode("errorDecision", func(c *Context) (bool, error) {
		return false, errors.New("decision error")
	})
	writeAction, _ := NewActionNode("writeAction", func(c *Context) (bool, error) {
		c.Store("write_action", "done")
		return true, nil
	})
	writeAnotherAction, _ := NewActionNode("writeAnotherAction", func(c *Context) (bool, error) {
		c.Store("write_another_action", "done")
		return true, nil
	})
	readAction, _ := NewActionNode("readAction", func(c *Context) (bool, error) {
		v, ok := c.Read("write_action")
		if !ok {
			return false, nil
		}
		c.Store("read_action", fmt.Sprintf("the content of write_action is %v", v))
		return true, nil
	})
	deleteAnotherAction, _ := NewActionNode("deleteAnotherAction", func(c *Context) (bool, error) {
		c.Delete("write_another_action")
		return true, nil
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
		givenContextData    contextData
		expectedIsDone      bool
		expectedContextData contextData
		expectedReport      map[Node]ComputeState
	}{
		{
			name:                "Can compute empty validated system",
			expectedIsDone:      true,
			expectedContextData: contextData{},
			expectedReport:      map[Node]ComputeState{},
		},
		{
			name: "Can compute one action node system",
			givenNodes: []Node{
				writeAction,
			},
			expectedIsDone: true,
			expectedContextData: contextData{
				"write_action": "done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction: ComputeStateContinue(),
			},
		},
		{
			name: "Can compute one stop action node system",
			givenNodes: []Node{
				readAction,
			},
			expectedIsDone:      true,
			expectedContextData: contextData{},
			expectedReport: map[Node]ComputeState{
				readAction: ComputeStateStop(),
			},
		},
		{
			name: "Can compute 2 action nodes system",
			givenNodes: []Node{
				writeAction,
				writeAnotherAction,
			},
			expectedIsDone: true,
			expectedContextData: contextData{
				"write_action":         "done",
				"write_another_action": "done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction:        ComputeStateContinue(),
				writeAnotherAction: ComputeStateContinue(),
			},
		},
		{
			name: "Can compute 2 linked action system (ordered declaration)",
			givenNodes: []Node{
				writeAction,
				readAction,
			},
			givenLinks: []nodeLink{
				newLink(writeAction, readAction),
			},
			expectedIsDone: true,
			expectedContextData: contextData{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction: ComputeStateContinue(),
				readAction:  ComputeStateContinue(),
			},
		},
		{
			name: "Can compute 2 linked action system (unordered declaration)",
			givenNodes: []Node{
				readAction,
				writeAction,
			},
			givenLinks: []nodeLink{
				newLink(writeAction, readAction),
			},
			expectedIsDone: true,
			expectedContextData: contextData{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction: ComputeStateContinue(),
				readAction:  ComputeStateContinue(),
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
				newLink(writeAction, writeActionKeyIsPresent),
				newBranchLink(writeActionKeyIsPresent, readAction, true),
				newBranchLink(writeActionKeyIsPresent, deleteAnotherAction, false),
			},
			expectedIsDone: true,
			expectedContextData: contextData{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction:             ComputeStateContinue(),
				writeActionKeyIsPresent: ComputeStateContinueOnBranch(true),
				readAction:              ComputeStateContinue(),
				deleteAnotherAction:     ComputeStateSkip(),
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
				newLink(writeAnotherAction, writeActionKeyIsPresent),
				newBranchLink(writeActionKeyIsPresent, readAction, true),
				newBranchLink(writeActionKeyIsPresent, deleteAnotherAction, false),
			},
			expectedIsDone:      true,
			expectedContextData: contextData{},
			expectedReport: map[Node]ComputeState{
				writeAnotherAction:      ComputeStateContinue(),
				writeActionKeyIsPresent: ComputeStateContinueOnBranch(false),
				readAction:              ComputeStateSkip(),
				deleteAnotherAction:     ComputeStateContinue(),
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
				newLink(writeAction, writeActionKeyIsPresent),
				newBranchLink(writeActionKeyIsPresent, readAction, true),
				newBranchLink(writeActionKeyIsPresent, deleteAnotherAction, false),
			},
			expectedIsDone: true,
			expectedContextData: contextData{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction:             ComputeStateContinue(),
				writeActionKeyIsPresent: ComputeStateContinueOnBranch(true),
				readAction:              ComputeStateContinue(),
				deleteAnotherAction:     ComputeStateSkip(),
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
				newLink(writeAnotherAction, writeActionKeyIsPresent),
				newBranchLink(writeActionKeyIsPresent, readAction, true),
				newBranchLink(writeActionKeyIsPresent, deleteAnotherAction, false),
			},
			expectedIsDone:      true,
			expectedContextData: contextData{},
			expectedReport: map[Node]ComputeState{
				writeAnotherAction:      ComputeStateContinue(),
				writeActionKeyIsPresent: ComputeStateContinueOnBranch(false),
				readAction:              ComputeStateSkip(),
				deleteAnotherAction:     ComputeStateContinue(),
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
				newLink(writeAction, writeActionKeyIsPresent),
				newBranchLink(writeActionKeyIsPresent, errorAction, true),
				newLink(errorAction, readAction),
				newBranchLink(writeActionKeyIsPresent, deleteAnotherAction, false),
			},
			expectedIsDone: false,
			expectedContextData: contextData{
				"write_action": "done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction:             ComputeStateContinue(),
				writeActionKeyIsPresent: ComputeStateContinueOnBranch(true),
				errorAction:             ComputeStateAbort(errors.New("action error")),
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
				newLink(writeAnotherAction, writeActionKeyIsPresent),
				newBranchLink(writeActionKeyIsPresent, errorAction, false),
				newLink(errorAction, readAction),
			},
			expectedIsDone: false,
			expectedContextData: contextData{
				"write_another_action": "done",
			},
			expectedReport: map[Node]ComputeState{
				writeAnotherAction:      ComputeStateContinue(),
				writeActionKeyIsPresent: ComputeStateContinueOnBranch(false),
				errorAction:             ComputeStateAbort(errors.New("action error")),
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
				newLink(writeAction, errorDecision),
				newBranchLink(errorDecision, readAction, true),
				newBranchLink(errorDecision, deleteAnotherAction, false),
			},
			expectedIsDone: false,
			expectedContextData: contextData{
				"write_action": "done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction:   ComputeStateContinue(),
				errorDecision: ComputeStateAbort(errors.New("decision error")),
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
				newLink(writeAction, deleteAnotherAction),
				newLink(writeAction, errorAction),
				newLink(writeAction, readAction),
			},
			expectedIsDone: false,
			expectedContextData: contextData{
				"write_action": "done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction:         ComputeStateContinue(),
				deleteAnotherAction: ComputeStateContinue(),
				errorAction:         ComputeStateAbort(errors.New("action error")),
			},
		},
		{
			name: "Can compute a node system with join links",
			givenNodes: []Node{
				writeAction,
				writeAnotherAction,
				readAction,
				deleteAnotherAction,
			},
			givenNodesJoinModes: map[Node]JoinMode{
				readAction: JoinModeAnd,
			},
			givenLinks: []nodeLink{
				newLink(writeAction, readAction),
				newLink(writeAnotherAction, readAction),
				newLink(readAction, deleteAnotherAction),
			},
			expectedIsDone: true,
			expectedContextData: contextData{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction:         ComputeStateContinue(),
				writeAnotherAction:  ComputeStateContinue(),
				readAction:          ComputeStateContinue(),
				deleteAnotherAction: ComputeStateContinue(),
			},
		},
		{
			name: "Can compute a node system with partial join links",
			givenNodes: []Node{
				writeAction,
				writeActionKeyIsPresent,
				writeAnotherAction,
				readAction,
				deleteAnotherAction,
			},
			givenNodesJoinModes: map[Node]JoinMode{
				readAction: JoinModeAnd,
			},
			givenLinks: []nodeLink{
				newLink(writeAction, readAction),
				newBranchLink(writeActionKeyIsPresent, readAction, false),
				newLink(writeAnotherAction, readAction),
				newLink(readAction, deleteAnotherAction),
			},
			expectedIsDone: true,
			expectedContextData: contextData{
				"write_action":         "done",
				"write_another_action": "done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction:             ComputeStateContinue(),
				writeActionKeyIsPresent: ComputeStateContinueOnBranch(true),
				writeAnotherAction:      ComputeStateContinue(),
			},
		},
		{
			name: "Can compute a node system with merge links",
			givenNodes: []Node{
				writeAction,
				writeActionKeyIsPresent,
				readAction,
				deleteAnotherAction,
			},
			givenNodesJoinModes: map[Node]JoinMode{
				readAction: JoinModeOr,
			},
			givenLinks: []nodeLink{
				newLink(writeAction, readAction),
				newBranchLink(writeActionKeyIsPresent, readAction, true),
				newLink(readAction, deleteAnotherAction),
			},
			expectedIsDone: true,
			expectedContextData: contextData{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction:             ComputeStateContinue(),
				writeActionKeyIsPresent: ComputeStateContinueOnBranch(true),
				readAction:              ComputeStateContinue(),
				deleteAnotherAction:     ComputeStateContinue(),
			},
		},
		{
			name: "Can compute a node system with partial merge links",
			givenNodes: []Node{
				writeActionKeyIsPresent,
				writeAnotherActionKeyIsPresent,
				readAction,
				deleteAnotherAction,
			},
			givenNodesJoinModes: map[Node]JoinMode{
				readAction: JoinModeOr,
			},
			givenLinks: []nodeLink{
				newBranchLink(writeActionKeyIsPresent, readAction, true),
				newBranchLink(writeAnotherActionKeyIsPresent, readAction, true),
				newLink(readAction, deleteAnotherAction),
			},
			expectedIsDone:      true,
			expectedContextData: contextData{},
			expectedReport: map[Node]ComputeState{
				writeActionKeyIsPresent:        ComputeStateContinueOnBranch(false),
				writeAnotherActionKeyIsPresent: ComputeStateContinueOnBranch(false),
				readAction:                     ComputeStateSkip(),
				deleteAnotherAction:            ComputeStateSkip(),
			},
		},
		{
			name: "Can compute a node system with merge links who generate error",
			givenNodes: []Node{
				writeAction,
				errorAction,
				readAction,
				deleteAnotherAction,
			},
			givenNodesJoinModes: map[Node]JoinMode{
				readAction: JoinModeOr,
			},
			givenLinks: []nodeLink{
				newLink(writeAction, readAction),
				newLink(errorAction, readAction),
				newLink(readAction, deleteAnotherAction),
			},
			expectedIsDone: false,
			expectedContextData: contextData{
				"write_action": "done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction: ComputeStateContinue(),
				errorAction: ComputeStateAbort(errors.New("action error")),
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			system := NewNodeSystem()
			for _, node := range testCase.givenNodes {
				system.AddNode(node)
			}
			for node, mode := range testCase.givenNodesJoinModes {
				system.AddNodeJoinMode(node, mode)
			}
			for _, link := range testCase.givenLinks {
				if link.branch == nil {
					system.AddLink(link.from, link.to)
				} else {
					system.AddBranchLink(link.from, link.to, *link.branch)
				}
			}
			_, errs := system.Validate()
			if errs != nil {
				t.Errorf("validation errors - %+v\n", errs)
			}
			system.activate()

			c, err := NewComputation(system, setupContext(testCase.givenContextData))
			if err != nil {
				t.Errorf("can't compute: %+v\n", err)
				t.FailNow()
			}

			c.Compute()
			if !cmp.Equal(c.isDone(), testCase.expectedIsDone) {
				t.Errorf("computation is done - got: %+v, want: %+v", c.isDone(), testCase.expectedIsDone)
			}
			if !cmp.Equal(c.context.data, testCase.expectedContextData) {
				t.Errorf("context data - got: %+v, want: %+v", c.context.data, testCase.expectedContextData)
			}
			if !cmp.Equal(c.report, testCase.expectedReport) {
				t.Errorf("report - got: %+v, want: %+v", c.report, testCase.expectedReport)
			}
		})
	}
}
