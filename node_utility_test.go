package flow

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type NodeTestCase struct {
	name                string
	givenContextData    contextData
	givenNode           Node
	expectedRunState    RunState
	expectedContextData contextData
}

func RunTestOnNode(t *testing.T, testCases []NodeTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testContext := setupContext(testCase.givenContextData)
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
