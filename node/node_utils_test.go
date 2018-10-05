package node

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/computestate"

	"github.com/rlespinasse/hoff/internal/utils"
)

type NodeTestCase struct {
	name                 string
	givenContextData     map[string]interface{}
	givenNode            Node
	expectedComputeState computestate.ComputeState
	expectedContextData  map[string]interface{}
}

func RunTestOnNode(t *testing.T, testCases []NodeTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.givenContextData == nil {
				testCase.givenContextData = make(map[string]interface{})
			}
			testContext := &Context{
				data: testCase.givenContextData,
			}
			testState := testCase.givenNode.Compute(testContext)

			if !cmp.Equal(testState, testCase.expectedComputeState, utils.ErrorComparator) {
				t.Errorf("context state - got: %+v, want: %+v", testState, testCase.expectedComputeState)
			}

			if testCase.expectedContextData != nil && !cmp.Equal(testContext.data, testCase.expectedContextData) {
				t.Errorf("context data - got: %+v, want: %+v", testContext.data, testCase.expectedContextData)
			}
		})
	}
}
