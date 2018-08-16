package flow

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_NewContext(t *testing.T) {
	c := NewContext()
	emptyData := contextData{}

	if !cmp.Equal(c.data, emptyData) {
		t.Errorf("context data - got: %#v, want: %#v", c.data, emptyData)
	}
}

func Test_Context_Read(t *testing.T) {
	testCases := []struct {
		name             string
		givenContextData contextData
		givenKey         string
		expectedValue    interface{}
		expectedError    error
	}{
		{
			"value",
			contextData{
				"key": "value",
			},
			"key",
			"value",
			nil,
		},
		{
			"error",
			contextData{},
			"key",
			nil,
			errors.New("unknown key: key"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := setupContext(testCase.givenContextData)
			value, err := c.Read(testCase.givenKey)

			if value != testCase.expectedValue {
				t.Errorf("value - got: %#v, want: %#v", value, testCase.expectedValue)
			}
			if err != nil && testCase.expectedError != nil {
				if err.Error() != testCase.expectedError.Error() {
					t.Errorf("error - got: %#v, want: %#v", err, testCase.expectedError)
				}
			} else if err != nil || testCase.expectedError != nil {
				t.Errorf("error - got: %#v, want: %#v", err, testCase.expectedError)
			}
		})
	}
}
