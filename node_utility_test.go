package flow

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type FlowNodeTestCase struct {
	name                string
	givenContextData    contextData
	givenNode           FlowNode
	expectedRunState    RunState
	expectedContextData contextData
}

func RunTestOnFlowNode(t *testing.T, testCases []FlowNodeTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testContext := setupFlowContext(testCase.givenContextData)
			testState := testCase.givenNode.Run(testContext)

			if !cmp.Equal(testState, testCase.expectedRunState, runStateEqualOpts) {
				t.Errorf("context data - got: %#v, want: %#v", testState, testCase.expectedRunState)
			}

			if !cmp.Equal(testContext.data, testCase.expectedContextData) {
				t.Errorf("context data - got: %#v, want: %#v", testContext.data, testCase.expectedContextData)
			}
		})
	}
}
