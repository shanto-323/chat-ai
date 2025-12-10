package config

import "testing"

type ConfigTestStruct struct {
	name     string
	testFunc func(t *testing.T)
}

var ConfigTests = []ConfigTestStruct{}

func TestConfig(t *testing.T) {
	for _, tc := range ConfigTests {
		t.Run(tc.name, tc.testFunc)
	}
}
