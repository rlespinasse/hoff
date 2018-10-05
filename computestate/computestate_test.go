package computestate

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rlespinasse/hoff/internal/utils"
	"github.com/rlespinasse/hoff/statetype"
)

func Test_ComputeState_Call(t *testing.T) {
	testCases := []struct {
		name                  string
		givenComputeStateCall func() ComputeState
		expectedState         statetype.StateType
		expectedNodeBranch    *bool
		expectedError         error
		expectedString        string
	}{
		{
			name:                  "Should generate a continue state",
			givenComputeStateCall: func() ComputeState { return Continue() },
			expectedState:         statetype.ContinueState,
			expectedString:        "'Continue'",
		},
		{
			name:                  "Should generate a continue state on branch 'true'",
			givenComputeStateCall: func() ComputeState { return ContinueOnBranch(true) },
			expectedState:         statetype.ContinueState,
			expectedNodeBranch:    utils.BoolPointer(true),
			expectedString:        "'Continue on true'",
		},
		{
			name:                  "Should generate a stop state",
			givenComputeStateCall: func() ComputeState { return Stop() },
			expectedState:         statetype.StopState,
			expectedString:        "'Stop'",
		},
		{
			name:                  "Should generate a skip state",
			givenComputeStateCall: func() ComputeState { return Skip() },
			expectedState:         statetype.SkipState,
			expectedString:        "'Skip'",
		},
		{
			name:                  "Should generate a abort state",
			givenComputeStateCall: func() ComputeState { return Abort(errors.New("error")) },
			expectedState:         statetype.AbortState,
			expectedError:         errors.New("error"),
			expectedString:        "'Abort on error'",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			computeState := testCase.givenComputeStateCall()
			if computeState.Value != testCase.expectedState {
				t.Errorf("state - got: %+v, want: %+v", computeState.Value, testCase.expectedState)
			}
			if !cmp.Equal(computeState.Branch, testCase.expectedNodeBranch) {
				t.Errorf("branch - got: %+v, want: %+v", computeState.Branch, testCase.expectedNodeBranch)
			}
			if !cmp.Equal(computeState.Error, testCase.expectedError, utils.ErrorComparator) {
				t.Errorf("error - got: %+v, want: %+v", computeState.Error, testCase.expectedError)
			}
			if computeState.String() != testCase.expectedString {
				t.Errorf("string - got: %+v, want: %+v", computeState.String(), testCase.expectedString)
			}
		})
	}
}
