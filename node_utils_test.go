package hoff

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

type SomeNode struct{}

func (n *SomeNode) Compute(c *Context) ComputeState {
	c.Store("message", "SomeNode is passing")
	return NewContinueComputeState()
}

func (n *SomeNode) DecideCapability() bool {
	return false
}

type AnotherNode struct{}

func (n *AnotherNode) Compute(c *Context) ComputeState {
	c.Store("message", "AnotherNode is passing")
	return NewContinueComputeState()
}

func (n *AnotherNode) DecideCapability() bool {
	return false
}

type NodeTestCase struct {
	name                 string
	givenContextData     map[string]interface{}
	givenNode            Node
	expectedComputeState ComputeState
	expectedContextData  map[string]interface{}
}

func RunTestOnNode(t *testing.T, testCases []NodeTestCase) {
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			if testCase.givenContextData == nil {
				testCase.givenContextData = make(map[string]interface{})
			}
			testContext := &Context{
				Data: testCase.givenContextData,
			}
			testState := testCase.givenNode.Compute(testContext)

			if !cmp.Equal(testState, testCase.expectedComputeState, errorComparator) {
				t.Errorf("context state - got: %+v, want: %+v", testState, testCase.expectedComputeState)
			}

			if testCase.expectedContextData != nil && !cmp.Equal(testContext.Data, testCase.expectedContextData) {
				t.Errorf("context data - got: %+v, want: %+v", testContext.Data, testCase.expectedContextData)
			}
		})
	}
}
