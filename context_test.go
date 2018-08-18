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
		t.Errorf("context data - got: %+v, want: %+v", c.data, emptyData)
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
			name:             "value",
			givenContextData: contextData{"key": "value"},
			givenKey:         "key",
			expectedValue:    "value",
		},
		{
			name:          "error",
			givenKey:      "key",
			expectedError: errors.New("unknown key: key"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := setupContext(testCase.givenContextData)
			value, err := c.Read(testCase.givenKey)

			if value != testCase.expectedValue {
				t.Errorf("value - got: %+v, want: %+v", value, testCase.expectedValue)
			}
			if !cmp.Equal(err, testCase.expectedError, errorEqualOpts) {
				t.Errorf("error - got: %+v, want: %+v", err, testCase.expectedError)
			}
		})
	}
}
			}
		})
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
			name:             "Can read present key",
			givenContextData: contextData{"key": "value"},
			givenKey:         "key",
			expectedValue:    "value",
		},
		{
			name:          "Can't read unknown key",
			givenKey:      "key",
			expectedError: errors.New("unknown key: key"),
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := setupContext(testCase.givenContextData)
			value, err := c.Read(testCase.givenKey)

			if value != testCase.expectedValue {
				t.Errorf("value - got: %+v, want: %+v", value, testCase.expectedValue)
			}
			if !cmp.Equal(err, testCase.expectedError, errorEqualOpts) {
				t.Errorf("error - got: %+v, want: %+v", err, testCase.expectedError)
			}
		})
	}
}

func Test_Context_Delete(t *testing.T) {
	testCases := []struct {
		name                string
		givenContextData    contextData
		givenKey            string
		expectedContextData contextData
	}{
		{
			name:                "Can delete a key",
			givenContextData:    contextData{"key": "value"},
			givenKey:            "key",
			expectedContextData: contextData{},
		},
		{
			name: "Can delete a present key without deleting other keys",
			givenContextData: contextData{
				"key":         "value",
				"another_key": "another_value",
			},
			givenKey: "another_key",
			expectedContextData: contextData{
				"key": "value",
			},
		},
		{
			name: "Can delete a missing key",
			givenContextData: contextData{
				"key": "value",
			},
			givenKey: "missing_key",
			expectedContextData: contextData{
				"key": "value",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := setupContext(testCase.givenContextData)
			c.Delete(testCase.givenKey)

			if !cmp.Equal(c.data, testCase.expectedContextData) {
				t.Errorf("got: %+v, want: %+v", c.data, testCase.expectedContextData)
			}
		})
	}
}

func Test_Context_Have(t *testing.T) {
	testCases := []struct {
		name             string
		givenContextData contextData
		givenKey         string
		expectedResult   bool
	}{
		{
			name:             "Can check if context have present key",
			givenContextData: contextData{"key": "value"},
			givenKey:         "key",
			expectedResult:   true,
		},
		{
			name:             "Can check if context have missing key",
			givenContextData: contextData{"key": "value"},
			givenKey:         "missing_key",
			expectedResult:   false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := setupContext(testCase.givenContextData)
			result := c.HaveKey(testCase.givenKey)

			if result != testCase.expectedResult {
				t.Errorf("got: %+v, want: %+v", result, testCase.expectedResult)
			}
		})
	}
}
