package main

import (
	"testing"
)

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("got an unexpected error, %v", err)
	}
}

func assertErrorMessage(t *testing.T, error error, want string) {
	t.Helper()
	errorMessage := error.Error()
	if errorMessage != want {
		t.Fatalf("error unexpected error message %v, wanted %v", errorMessage, want)
	}
}

func assertStringOutput(t *testing.T, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %v but wanted %v", got, want)
	}
}
