package flow

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type NodeTestCase struct {
	name                 string
	givenContextData     contextData
	givenNode            Node
	expectedComputeState ComputeState
	expectedContextData  contextData
}

func RunTestOnNode(t *testing.T, testCases []NodeTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testContext := setupContext(testCase.givenContextData)
			testState := testCase.givenNode.Compute(testContext)

			if !cmp.Equal(testState, testCase.expectedComputeState, computeStateEqualOpts) {
				t.Errorf("context state - got: %+v, want: %+v", testState, testCase.expectedComputeState)
			}

			if testCase.expectedContextData != nil && !cmp.Equal(testContext.data, testCase.expectedContextData) {
				t.Errorf("context data - got: %+v, want: %+v", testContext.data, testCase.expectedContextData)
			}
		})
	}
}
