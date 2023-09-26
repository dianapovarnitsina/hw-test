package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReadDir(t *testing.T) {
	dir, err := os.MkdirTemp("", "testenv")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(dir) // Удаляем временную директорию после завершения тестов

	err = os.WriteFile(filepath.Join(dir, "file1"), []byte("some_value1"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	err = os.WriteFile(filepath.Join(dir, "file2"), []byte{}, 0o644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	err = os.WriteFile(filepath.Join(dir, "file3"), []byte("some_value3"), 0o644)
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}

	expectedEnv := Environment{
		"file1": {"some_value1", false},
		"file2": {"", true},
		"file3": {"some_value3", false},
	}

	env, err := ReadDir(dir)
	if err != nil {
		t.Fatalf("Failed to read directory: %v", err)
	}

	for key, expectedValue := range expectedEnv {
		if envValue, ok := env[key]; ok {
			if envValue.NeedRemove != expectedValue.NeedRemove {
				t.Fatalf("Expected value for key '%s' to be { %v}, got { %v}", key, expectedValue.NeedRemove, envValue.NeedRemove)
			}
		} else {
			t.Fatalf("Expected key '%s' not found in environment", key)
		}
	}
}
