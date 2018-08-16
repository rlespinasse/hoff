package flowengine

import (
	"reflect"
	"testing"
)

type FlowNodeTestCase struct {
	name                string
	givenContextData    contextData
	givenNode           FlowNode
	expectedContextData contextData
}

func RunTestOnFlowNode(t *testing.T, testCases []FlowNodeTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			testContext := setupFlowContext(testCase.givenContextData)
			testCase.givenNode.Run(testContext)

			if !reflect.DeepEqual(testCase.expectedContextData, testContext.data) {
				t.Errorf("context data - got: %#v, want: %#v", testContext.data, testCase.expectedContextData)
			}
		})
	}
}
