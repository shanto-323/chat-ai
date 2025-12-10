package logger

import "testing"

type LoggerTestStruct struct {
	name     string
	testFunc func(t *testing.T)
}

var ConfigTests = []LoggerTestStruct{}

func TestConfig(t *testing.T) {
	for _, tc := range ConfigTests {
		t.Run(tc.name, tc.testFunc)
	}
}

