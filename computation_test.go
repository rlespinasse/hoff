package flow

import (
	"errors"
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

			if !cmp.Equal(c, testCase.expectedComputation, computationEqualOpts) {
				t.Errorf("computation - got: %+v, want: %+v", c, testCase.expectedComputation)
			}
			if !cmp.Equal(err, testCase.expectedError, errorEqualOpts) {
				t.Errorf("error - got: %+v, want: %+v", err, testCase.expectedError)
			}
		})
	}
}

var computationEqualOpts = cmp.Comparer(func(x, y *computation) bool {
	return (x == nil && y == nil) || (x != nil && y != nil && *x == *y)
})
