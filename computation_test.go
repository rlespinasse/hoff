package flow

import (
	"errors"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_NewComputation(t *testing.T) {
	var validatedSystem = NewNodeSystem()
	validatedSystem.validity = true

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
			givenSystem:         emptySystem,
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

			if !cmp.Equal(c, testCase.expectedComputation, pointerOfComputationEqualOpts) {
				t.Errorf("computation - got: %+v, want: %+v", c, testCase.expectedComputation)
			}
			if !cmp.Equal(err, testCase.expectedError, errorEqualOpts) {
				t.Errorf("error - got: %+v, want: %+v", err, testCase.expectedError)
			}
		})
	}
}

func Test_Computation_Compute(t *testing.T) {
	stopAction := NewActionNode(func(c *Context) (bool, error) {
		return false, nil
	})
	errorAction := NewActionNode(func(c *Context) (bool, error) {
		return false, errors.New("action error")
	})
	errorDecision := NewDecisionNode(func(c *Context) (bool, error) {
		return false, errors.New("decision error")
	})
	writeAction := NewActionNode(func(c *Context) (bool, error) {
		c.Store("write_action", "done")
		return true, nil
	})
	writeAnotherAction := NewActionNode(func(c *Context) (bool, error) {
		c.Store("write_another_action", "done")
		return true, nil
	})
	readAction := NewActionNode(func(c *Context) (bool, error) {
		v, err := c.Read("write_action")
		if err != nil {
			return false, err
		}
		c.Store("read_action", fmt.Sprintf("the content of write_action is %v", v))
		return true, nil
	})
	deleteAnotherAction := NewActionNode(func(c *Context) (bool, error) {
		c.Delete("write_another_action")
		return true, nil
	})
	writeActionKeyIsPresent := NewDecisionNode(func(c *Context) (bool, error) {
		return c.HaveKey("write_action"), nil
	})
	testCases := []struct {
		name                string
		givenNodes          []Node
		givenBranchLinks    []NodeBranchLink
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
				writeAction: ComputeStatePass(),
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
				writeAction:        ComputeStatePass(),
				writeAnotherAction: ComputeStatePass(),
			},
		},
		{
			name: "Can compute 2 linked action system (ordered declaration)",
			givenNodes: []Node{
				writeAction,
				readAction,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: writeAction,
					To:   readAction,
				},
			},
			expectedIsDone: true,
			expectedContextData: contextData{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction: ComputeStatePass(),
				readAction:  ComputeStatePass(),
			},
		},
		{
			name: "Can compute 2 linked action system (unordered declaration)",
			givenNodes: []Node{
				readAction,
				writeAction,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: writeAction,
					To:   readAction,
				},
			},
			expectedIsDone: true,
			expectedContextData: contextData{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction: ComputeStatePass(),
				readAction:  ComputeStatePass(),
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
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: writeAction,
					To:   writeActionKeyIsPresent,
				},
				NodeBranchLink{
					From:   writeActionKeyIsPresent,
					To:     readAction,
					Branch: ptrOfString("true"),
				},
				NodeBranchLink{
					From:   writeActionKeyIsPresent,
					To:     deleteAnotherAction,
					Branch: ptrOfString("false"),
				},
			},
			expectedIsDone: true,
			expectedContextData: contextData{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction:             ComputeStatePass(),
				writeActionKeyIsPresent: ComputeStateBranchPass("true"),
				readAction:              ComputeStatePass(),
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
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: writeAnotherAction,
					To:   writeActionKeyIsPresent,
				},
				NodeBranchLink{
					From:   writeActionKeyIsPresent,
					To:     readAction,
					Branch: ptrOfString("true"),
				},
				NodeBranchLink{
					From:   writeActionKeyIsPresent,
					To:     deleteAnotherAction,
					Branch: ptrOfString("false"),
				},
			},
			expectedIsDone:      true,
			expectedContextData: contextData{},
			expectedReport: map[Node]ComputeState{
				writeAnotherAction:      ComputeStatePass(),
				writeActionKeyIsPresent: ComputeStateBranchPass("false"),
				deleteAnotherAction:     ComputeStatePass(),
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
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: writeAction,
					To:   writeActionKeyIsPresent,
				},
				NodeBranchLink{
					From:   writeActionKeyIsPresent,
					To:     readAction,
					Branch: ptrOfString("true"),
				},
				NodeBranchLink{
					From:   writeActionKeyIsPresent,
					To:     deleteAnotherAction,
					Branch: ptrOfString("false"),
				},
			},
			expectedIsDone: true,
			expectedContextData: contextData{
				"write_action": "done",
				"read_action":  "the content of write_action is done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction:             ComputeStatePass(),
				writeActionKeyIsPresent: ComputeStateBranchPass("true"),
				readAction:              ComputeStatePass(),
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
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: writeAnotherAction,
					To:   writeActionKeyIsPresent,
				},
				NodeBranchLink{
					From:   writeActionKeyIsPresent,
					To:     readAction,
					Branch: ptrOfString("true"),
				},
				NodeBranchLink{
					From:   writeActionKeyIsPresent,
					To:     deleteAnotherAction,
					Branch: ptrOfString("false"),
				},
			},
			expectedIsDone:      true,
			expectedContextData: contextData{},
			expectedReport: map[Node]ComputeState{
				writeAnotherAction:      ComputeStatePass(),
				writeActionKeyIsPresent: ComputeStateBranchPass("false"),
				deleteAnotherAction:     ComputeStatePass(),
			},
		},
		{
			name:                "Can compute one stopping action node system",
			givenNodes:          []Node{stopAction},
			expectedIsDone:      true,
			expectedContextData: contextData{},
			expectedReport:      map[Node]ComputeState{stopAction: ComputeStateStop()},
		},
		{
			name: "Can compute a node system with one stopping action node",
			givenNodes: []Node{
				writeAction,
				writeActionKeyIsPresent,
				stopAction,
				readAction,
				deleteAnotherAction,
			},
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: writeAction,
					To:   writeActionKeyIsPresent,
				},
				NodeBranchLink{
					From:   writeActionKeyIsPresent,
					To:     stopAction,
					Branch: ptrOfString("true"),
				},
				NodeBranchLink{
					From: stopAction,
					To:   readAction,
				},
				NodeBranchLink{
					From:   writeActionKeyIsPresent,
					To:     deleteAnotherAction,
					Branch: ptrOfString("false"),
				},
			},
			expectedIsDone: true,
			expectedContextData: contextData{
				"write_action": "done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction:             ComputeStatePass(),
				writeActionKeyIsPresent: ComputeStateBranchPass("true"),
				stopAction:              ComputeStateStop(),
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
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: writeAction,
					To:   writeActionKeyIsPresent,
				},
				NodeBranchLink{
					From:   writeActionKeyIsPresent,
					To:     errorAction,
					Branch: ptrOfString("true"),
				},
				NodeBranchLink{
					From: errorAction,
					To:   readAction,
				},
				NodeBranchLink{
					From:   writeActionKeyIsPresent,
					To:     deleteAnotherAction,
					Branch: ptrOfString("false"),
				},
			},
			expectedIsDone: false,
			expectedContextData: contextData{
				"write_action": "done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction:             ComputeStatePass(),
				writeActionKeyIsPresent: ComputeStateBranchPass("true"),
				errorAction:             ComputeStateFail(errors.New("action error")),
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
			givenBranchLinks: []NodeBranchLink{
				NodeBranchLink{
					From: writeAction,
					To:   errorDecision,
				},
				NodeBranchLink{
					From:   errorDecision,
					To:     readAction,
					Branch: ptrOfString("true"),
				},
				NodeBranchLink{
					From:   errorDecision,
					To:     deleteAnotherAction,
					Branch: ptrOfString("false"),
				},
			},
			expectedIsDone: false,
			expectedContextData: contextData{
				"write_action": "done",
			},
			expectedReport: map[Node]ComputeState{
				writeAction:   ComputeStatePass(),
				errorDecision: ComputeStateFail(errors.New("decision error")),
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			system := NewNodeSystem()
			for _, node := range testCase.givenNodes {
				system.AddNode(node)
			}
			for _, link := range testCase.givenBranchLinks {
				system.AddBranchLink(link)
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
