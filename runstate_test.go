package flow

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_State_String(t *testing.T) {
	testCases := []struct {
		name            string
		givenStateValue State
		expectedString  string
	}{
		{
			"Print pass",
			pass,
			"pass",
		},
		{
			"Print stop",
			stop,
			"stop",
		},
		{
			"Print Fail",
			fail,
			"fail",
		},
		{
			"Print Unknowned state",
			0,
			"unknown",
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

func Test_AvailableBranches(t *testing.T) {
	testCases := []struct {
		name                 string
		givenBranches        []string
		expectedNodeBranches []NodeBranch
	}{
		{
			"With branches",
			[]string{"branch"},
			[]NodeBranch{
				newNodeBranch("branch"),
			},
		},
		{
			"Without branches",
			[]string{},
			[]NodeBranch{},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			branches := AvailablesBranches(testCase.givenBranches...)
			if !cmp.Equal(branches, testCase.expectedNodeBranches) {
				t.Errorf("got: %#v, want: %#v", branches, testCase.expectedNodeBranches)
			}
		})
	}
}

func Test_RunState_Call(t *testing.T) {
	testCases := []struct {
		name               string
		givenRunStateCall  func() RunState
		expectedState      State
		expectedNodeBranch NodeBranch
		expectedError      error
	}{
		{
			"Pass",
			func() RunState {
				return RunStatePass()
			},
			pass,
			nil,
			nil,
		},
		{
			"BranchPass",
			func() RunState {
				return RunStateBranchPass("branch")
			},
			pass,
			newNodeBranch("branch"),
			nil,
		},
		{
			"Stop",
			func() RunState {
				return RunStateStop()
			},
			stop,
			nil,
			nil,
		},
		{
			"BranchStop",
			func() RunState {
				return RunStateBranchStop("branch")
			},
			stop,
			newNodeBranch("branch"),
			nil,
		},
		{
			"Fail",
			func() RunState {
				return RunStateFail(errors.New("error"))
			},
			fail,
			nil,
			errors.New("error"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			runState := testCase.givenRunStateCall()
			if runState.value != testCase.expectedState {
				t.Errorf("state - got: %#v, want: %#v", runState.value, testCase.expectedState)
			}
			if runState.branch != nil && testCase.expectedNodeBranch != nil {
				if *runState.branch != *testCase.expectedNodeBranch {
					t.Errorf("branch - got: %#v, want: %#v", runState.branch, testCase.expectedNodeBranch)
				}
			} else if runState.branch != nil || testCase.expectedNodeBranch != nil {
				t.Errorf("branch - got: %#v, want: %#v", runState.branch, testCase.expectedNodeBranch)
			}
			if runState.err != nil && testCase.expectedError != nil {
				if runState.err.Error() != testCase.expectedError.Error() {
					t.Errorf("err - got: %#v, want: %#v", runState.err, testCase.expectedError)
				}
			} else if runState.err != nil || testCase.expectedError != nil {
				t.Errorf("err - got: %#v, want: %#v", runState.err, testCase.expectedError)
			}
		})
	}
}
