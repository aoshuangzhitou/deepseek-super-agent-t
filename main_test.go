package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestGreeting(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{name: "default name", input: "", expected: "Hello, Human!"},
		{name: "explicit name", input: "Ada", expected: "Hello, Ada!"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := greeting(tt.input); got != tt.expected {
				t.Fatalf("greeting(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestRun(t *testing.T) {
	var out bytes.Buffer
	if err := run(&out, []string{"-name", "Grace"}); err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if got, want := out.String(), "Hello, Grace!\n"; got != want {
		t.Fatalf("run output = %q, want %q", got, want)
	}
}

func TestRunDefault(t *testing.T) {
	var out bytes.Buffer
	if err := run(&out, nil); err != nil {
		t.Fatalf("run returned error: %v", err)
	}
	if !strings.Contains(out.String(), "Hello, Human!") {
		t.Fatalf("run output = %q, want it to contain %q", out.String(), "Hello, Human!")
	}
}
