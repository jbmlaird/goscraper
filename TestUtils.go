package main

import "testing"

func assertNoError(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("got an unexpected error, %v", err)
	}
}

func assertErrorMessage(t *testing.T, got string, want string) {
	t.Helper()
	if got != want {
		t.Fatalf("got unexpected error message %v, wanted %v", got, want)
	}
}

func assertOutput(t *testing.T, got, want bool) {
	t.Helper()
	if got != want {
		t.Errorf("got %v but wanted %v", got, want)
	}
}
