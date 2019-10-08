package hoff

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_NewContext(t *testing.T) {
	c := NewContext(map[string]interface{}{
		"key": "value",
	})
	expectedData := map[string]interface{}{
		"key": "value",
	}

	if !cmp.Equal(c.Data, expectedData) {
		t.Errorf("context data - got: %+v, want: %+v", c.Data, expectedData)
	}
}

func Test_NewContextWithoutData(t *testing.T) {
	c := NewContextWithoutData()
	emptyData := map[string]interface{}{}

	if !cmp.Equal(c.Data, emptyData) {
		t.Errorf("context data - got: %+v, want: %+v", c.Data, emptyData)
	}
}

func Test_Context_Equal(t *testing.T) {
	c1 := NewContext(map[string]interface{}{
		"key": "value",
	})
	c2 := NewContext(map[string]interface{}{
		"key": "value",
	})

	result := cmp.Equal(c1, c2)
	if !result {
		t.Errorf("got: %+v, want: %+v", result, true)
	}
}

func Test_Context_NotEqual(t *testing.T) {
	c1 := NewContext(map[string]interface{}{
		"key": "value",
	})
	c2 := NewContext(map[string]interface{}{
		"other_key": "value",
	})

	result := cmp.Equal(c1, c2)
	if result {
		t.Errorf("got: %+v, want: %+v", result, false)
	}
}

func Test_Context_Store(t *testing.T) {
	testCases := []struct {
		name                string
		givenKey            string
		givenValue          interface{}
		expectedContextData map[string]interface{}
	}{
		{
			name:                "Can store key:value",
			givenKey:            "key",
			givenValue:          "value",
			expectedContextData: map[string]interface{}{"key": "value"},
		},
		{
			name:                "Can store key:nil",
			givenKey:            "key",
			givenValue:          nil,
			expectedContextData: map[string]interface{}{"key": nil},
		},
		{
			name:                "Can store key:interface",
			givenKey:            "key",
			givenValue:          map[string]string{"map_key": "value"},
			expectedContextData: map[string]interface{}{"key": map[string]string{"map_key": "value"}},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := NewContextWithoutData()
			c.Store(testCase.givenKey, testCase.givenValue)

			if !cmp.Equal(c.Data, testCase.expectedContextData) {
				t.Errorf("got: %+v, want: %+v", c.Data, testCase.expectedContextData)
			}
		})
	}
}

func Test_Context_Read(t *testing.T) {
	testCases := []struct {
		name             string
		givenContextData map[string]interface{}
		givenKey         string
		expectedValue    interface{}
		expectedBool     bool
	}{
		{
			name:             "Can read present key",
			givenContextData: map[string]interface{}{"key": "value"},
			givenKey:         "key",
			expectedValue:    "value",
			expectedBool:     true,
		},
		{
			name:         "Can't read unknown key",
			givenKey:     "key",
			expectedBool: false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := &Context{
				Data: testCase.givenContextData,
			}
			value, ok := c.Read(testCase.givenKey)

			if value != testCase.expectedValue {
				t.Errorf("value - got: %+v, want: %+v", value, testCase.expectedValue)
			}
			if testCase.expectedBool != ok {
				t.Errorf("bool - got: %+v, want: %+v", ok, testCase.expectedBool)
			}
		})
	}
}

func Test_Context_Delete(t *testing.T) {
	testCases := []struct {
		name                string
		givenContextData    map[string]interface{}
		givenKey            string
		expectedContextData map[string]interface{}
	}{
		{
			name:                "Can delete a key",
			givenContextData:    map[string]interface{}{"key": "value"},
			givenKey:            "key",
			expectedContextData: map[string]interface{}{},
		},
		{
			name: "Can delete a present key without deleting other keys",
			givenContextData: map[string]interface{}{
				"key":         "value",
				"another_key": "another_value",
			},
			givenKey: "another_key",
			expectedContextData: map[string]interface{}{
				"key": "value",
			},
		},
		{
			name: "Can delete a missing key",
			givenContextData: map[string]interface{}{
				"key": "value",
			},
			givenKey: "missing_key",
			expectedContextData: map[string]interface{}{
				"key": "value",
			},
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := &Context{
				Data: testCase.givenContextData,
			}
			c.Delete(testCase.givenKey)

			if !cmp.Equal(c.Data, testCase.expectedContextData) {
				t.Errorf("got: %+v, want: %+v", c.Data, testCase.expectedContextData)
			}
		})
	}
}

func Test_Context_Have(t *testing.T) {
	testCases := []struct {
		name             string
		givenContextData map[string]interface{}
		givenKey         string
		expectedResult   bool
	}{
		{
			name:             "Can check if context have present key",
			givenContextData: map[string]interface{}{"key": "value"},
			givenKey:         "key",
			expectedResult:   true,
		},
		{
			name:             "Can check if context have missing key",
			givenContextData: map[string]interface{}{"key": "value"},
			givenKey:         "missing_key",
			expectedResult:   false,
		},
	}
	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			c := &Context{
				Data: testCase.givenContextData,
			}
			result := c.HaveKey(testCase.givenKey)

			if result != testCase.expectedResult {
				t.Errorf("got: %+v, want: %+v", result, testCase.expectedResult)
			}
		})
	}
}
