package main

import "testing"

type testEnvironment struct {
	inputEnv    Environment
	expectedEnv Environment
	expectedRC  int
}

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name    string
		testEnv testEnvironment
	}{
		{
			name: "Test with empty command",
			testEnv: testEnvironment{
				inputEnv: Environment{
					"TEST_VAR": {"test_value", false},
				},
				expectedEnv: Environment{
					"TEST_VAR": {"test_value", false},
				},
				expectedRC: 0,
			},
		},
		{
			name: "Test with valid command",
			testEnv: testEnvironment{
				inputEnv: Environment{
					"TEST_VAR": {"test_value", false},
				},
				expectedEnv: Environment{
					"TEST_VAR": {"test_value", false},
				},
				expectedRC: 0,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := []string{"echo", "hello world"}
			rc := RunCmd(cmd, tt.testEnv.inputEnv)

			if rc != tt.testEnv.expectedRC {
				t.Errorf("Expected return code %d, got %d", tt.testEnv.expectedRC, rc)
			}

			for key, expectedValue := range tt.testEnv.expectedEnv {
				envValue, ok := tt.testEnv.inputEnv[key]
				if !ok {
					t.Errorf("Expected environment key '%s' not found", key)
					continue
				}

				if envValue != expectedValue {
					t.Errorf("Expected value for key '%s' to be %v, got %v", key, expectedValue, envValue)
				}
			}
		})
	}
}
