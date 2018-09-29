package hoff

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_ComputeState_Call(t *testing.T) {
	testCases := []struct {
		name                  string
		givenComputeStateCall func() ComputeState
		expectedState         State
		expectedNodeBranch    *bool
		expectedError         error
		expectedString        string
	}{
		{
			name:                  "Should generate a continue state",
			givenComputeStateCall: func() ComputeState { return ComputeStateContinue() },
			expectedState:         Continue,
			expectedString:        "'Continue'",
		},
		{
			name:                  "Should generate a continue state on branch 'true'",
			givenComputeStateCall: func() ComputeState { return ComputeStateContinueOnBranch(true) },
			expectedState:         Continue,
			expectedNodeBranch:    boolPointer(true),
			expectedString:        "'Continue on true'",
		},
		{
			name:                  "Should generate a stop state",
			givenComputeStateCall: func() ComputeState { return ComputeStateStop() },
			expectedState:         Stop,
			expectedString:        "'Stop'",
		},
		{
			name:                  "Should generate a skip state",
			givenComputeStateCall: func() ComputeState { return ComputeStateSkip() },
			expectedState:         Skip,
			expectedString:        "'Skip'",
		},
		{
			name:                  "Should generate a abort state",
			givenComputeStateCall: func() ComputeState { return ComputeStateAbort(errors.New("error")) },
			expectedState:         Abort,
			expectedError:         errors.New("error"),
			expectedString:        "'Abort on error'",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			computeState := testCase.givenComputeStateCall()
			if computeState.value != testCase.expectedState {
				t.Errorf("state - got: %+v, want: %+v", computeState.value, testCase.expectedState)
			}
			if !cmp.Equal(computeState.branch, testCase.expectedNodeBranch) {
				t.Errorf("branch - got: %+v, want: %+v", computeState.branch, testCase.expectedNodeBranch)
			}
			if !cmp.Equal(computeState.err, testCase.expectedError, equalOptionForError) {
				t.Errorf("error - got: %+v, want: %+v", computeState.err, testCase.expectedError)
			}
			if computeState.String() != testCase.expectedString {
				t.Errorf("string - got: %+v, want: %+v", computeState.String(), testCase.expectedString)
			}
		})
	}
}
