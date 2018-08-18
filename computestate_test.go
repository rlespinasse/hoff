package flow

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
		expectedNodeBranch    *string
		expectedError         error
	}{
		{
			name: "Should generate a passing state",
			givenComputeStateCall: func() ComputeState { return ComputeStatePass() },
			expectedState:         pass,
		},
		{
			name: "Should generate a passing state on branch 'branch'",
			givenComputeStateCall: func() ComputeState { return ComputeStateBranchPass("branch") },
			expectedState:         pass,
			expectedNodeBranch:    ptrOfString("branch"),
		},
		{
			name: "Should generate a stopped state",
			givenComputeStateCall: func() ComputeState { return ComputeStateStop() },
			expectedState:         stop,
		},
		{
			name: "Should generate a stopped state on branch 'branch'",
			givenComputeStateCall: func() ComputeState { return ComputeStateBranchStop("branch") },
			expectedState:         stop,
			expectedNodeBranch:    ptrOfString("branch"),
		},
		{
			name: "Should generate a fail state",
			givenComputeStateCall: func() ComputeState { return ComputeStateStopOnError(errors.New("error")) },
			expectedState:         stop,
			expectedError:         errors.New("error"),
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
		})
	}
}
