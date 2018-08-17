package flow

import (
	"errors"
	"testing"
)

func Test_State_String(t *testing.T) {
	testCases := []struct {
		name            string
		givenStateValue State
		expectedString  string
	}{
		{
			name:            "Print pass",
			givenStateValue: pass,
			expectedString:  "pass",
		},
		{
			name:            "Print stop",
			givenStateValue: stop,
			expectedString:  "stop",
		},
		{
			name:            "Print Fail",
			givenStateValue: fail,
			expectedString:  "fail",
		},
		{
			name:            "Print Unknowned state",
			givenStateValue: 0,
			expectedString:  "unknown",
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			resultString := testCase.givenStateValue.String()
			if resultString != testCase.expectedString {
				t.Errorf("got: %v, want: %v", resultString, testCase.expectedString)
			}
		})
	}
}

func Test_ComputeState_Call(t *testing.T) {
	testCases := []struct {
		name                  string
		givenComputeStateCall func() ComputeState
		expectedState         State
		expectedNodeBranch    *string
		expectedError         error
	}{
		{
			"Pass",
			func() ComputeState {
				return ComputeStatePass()
			},
			pass,
			nil,
			nil,
		},
		{
			"BranchPass",
			func() ComputeState {
				return ComputeStateBranchPass("branch")
			},
			pass,
			ptrOfString("branch"),
			nil,
		},
		{
			"Stop",
			func() ComputeState {
				return ComputeStateStop()
			},
			stop,
			nil,
			nil,
		},
		{
			"BranchStop",
			func() ComputeState {
				return ComputeStateBranchStop("branch")
			},
			stop,
			ptrOfString("branch"),
			nil,
		},
		{
			"Fail",
			func() ComputeState {
				return ComputeStateFail(errors.New("error"))
			},
			fail,
			nil,
			errors.New("error"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			ComputeState := testCase.givenComputeStateCall()
			if ComputeState.value != testCase.expectedState {
				t.Errorf("state - got: %+v, want: %+v", ComputeState.value, testCase.expectedState)
			}
			if ComputeState.branch != nil && testCase.expectedNodeBranch != nil {
				if *ComputeState.branch != *testCase.expectedNodeBranch {
					t.Errorf("branch - got: %+v, want: %+v", ComputeState.branch, testCase.expectedNodeBranch)
				}
			} else if ComputeState.branch != nil || testCase.expectedNodeBranch != nil {
				t.Errorf("branch - got: %+v, want: %+v", ComputeState.branch, testCase.expectedNodeBranch)
			}
			if ComputeState.err != nil && testCase.expectedError != nil {
				if ComputeState.err.Error() != testCase.expectedError.Error() {
					t.Errorf("err - got: %+v, want: %+v", ComputeState.err, testCase.expectedError)
				}
			} else if ComputeState.err != nil || testCase.expectedError != nil {
				t.Errorf("err - got: %+v, want: %+v", ComputeState.err, testCase.expectedError)
			}
		})
	}
}
