package hoff

import (
	"errors"
	"fmt"
	"testing"

	"github.com/rlespinasse/hoff/computestate"
	"github.com/rlespinasse/hoff/internal/utils"

	"github.com/google/go-cmp/cmp"
)

func Test_Engine_ConfigureNodeSystem(t *testing.T) {
	activatedNodeSystem := NewNodeSystem()
	activatedNodeSystem.Activate()

	configuredEngine := &Engine{
		mode:   SequentialComputation,
		system: activatedNodeSystem,
	}

	testCases := []struct {
		name            string
		givenEngine     *Engine
		givenNodeSystem *NodeSystem
		expectedEngine  *Engine
		expectedError   error
	}{
		{
			name:            "Configure activated node system",
			givenEngine:     NewEngine(SequentialComputation),
			givenNodeSystem: activatedNodeSystem,
			expectedEngine:  configuredEngine,
		},
		{
			name:            "Configure unactivated node system",
			givenEngine:     NewEngine(SequentialComputation),
			givenNodeSystem: NewNodeSystem(),
			expectedEngine:  NewEngine(SequentialComputation),
			expectedError:   fmt.Errorf("node system need to be activated"),
		},
		{
			name:            "Try to reconfigure node system",
			givenEngine:     configuredEngine,
			givenNodeSystem: activatedNodeSystem,
			expectedEngine:  configuredEngine,
			expectedError:   fmt.Errorf("node system already configured"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			err := testCase.givenEngine.ConfigureNodeSystem(testCase.givenNodeSystem)

			if !cmp.Equal(testCase.givenEngine, testCase.expectedEngine, engineComparator) {
				t.Errorf("engine - got: %+v, want: %+v", testCase.givenEngine, testCase.expectedEngine)
			}
			if !cmp.Equal(err, testCase.expectedError, utils.ErrorComparator) {
				t.Errorf("error - got: %+v, want: %+v", err, testCase.expectedError)
			}
		})
	}
}

func Test_Engine_Compute(t *testing.T) {
	keyIsPresent, _ := NewDecisionNode("keyIsPresent", func(c *Context) (bool, error) {
		return c.HaveKey("key"), nil
	})
	stringAction, _ := NewActionNode("stringAction", func(c *Context) error {
		keyValue, _ := c.Read("key")
		c.Store("string", fmt.Sprintf("'%+v'", keyValue))
		return nil
	})
	throwedError := errors.New("missing 'key' in context")
	throwError, _ := NewActionNode("throwError", func(c *Context) error {
		return throwedError
	})

	ns := NewNodeSystem()
	ns.AddNode(keyIsPresent)
	ns.AddNode(stringAction)
	ns.AddNode(throwError)
	ns.AddLinkOnBranch(keyIsPresent, stringAction, true)
	ns.AddLinkOnBranch(keyIsPresent, throwError, false)
	ns.Activate()

	eng := NewEngine(SequentialComputation)
	eng.ConfigureNodeSystem(ns)

	testCases := []struct {
		name           string
		givenData      map[string]interface{}
		expectedResult ComputationResult
	}{
		{
			name:      "Compute with an generated error",
			givenData: make(map[string]interface{}),
			expectedResult: ComputationResult{
				Data:  make(map[string]interface{}),
				Error: throwedError,
				Report: map[Node]computestate.ComputeState{
					keyIsPresent: computestate.ContinueOnBranch(false),
					stringAction: computestate.Skip(),
					throwError:   computestate.Abort(throwedError),
				},
			},
		},
		{
			name: "Compute without error",
			givenData: map[string]interface{}{
				"key": []string{"Compute", "without", "error"},
			},
			expectedResult: ComputationResult{
				Data: map[string]interface{}{
					"key":    []string{"Compute", "without", "error"},
					"string": "'[Compute without error]'",
				},
				Report: map[Node]computestate.ComputeState{
					keyIsPresent: computestate.ContinueOnBranch(true),
					stringAction: computestate.Continue(),
					throwError:   computestate.Skip(),
				},
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result := eng.Compute(testCase.givenData)

			if !cmp.Equal(result, testCase.expectedResult, NodeComparator, utils.ErrorComparator) {
				t.Errorf("got: %+v, want: %+v", result, testCase.expectedResult)
			}
		})
	}
}

func Test_UnconfiguredEngine_Compute(t *testing.T) {
	eng := NewEngine(SequentialComputation)
	data := make(map[string]interface{})
	result := eng.Compute(data)

	expectedResult := ComputationResult{
		Data:  data,
		Error: errors.New("need a configured node system"),
	}

	if !cmp.Equal(result, expectedResult, NodeComparator, utils.ErrorComparator) {
		t.Errorf("got: %+v, want: %+v", result, expectedResult)
	}
}

var (
	engineComparator = cmp.Comparer(func(x, y Engine) bool {
		return x.mode == y.mode && ((x.system == nil && y.system == nil) || (x.system != nil && y.system != nil && cmp.Equal(x.system, y.system)))
	})
)
