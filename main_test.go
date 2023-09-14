package main

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestLsDir(t *testing.T) {
	var tests = []struct {
		path          string
		recursive     bool
		longFormat    bool
		showHidden    bool
		reverse       bool
		sortByModTime bool
		expected      string // Add an expected string output for each test case
	}{
		{".", false, false, false, false, false, "expected_output_here"},
		// Add more test cases here
	}

	for _, testCase := range tests {
		// Capture stdout
		old := os.Stdout
		r, w, _ := os.Pipe()
		os.Stdout = w

		lsDir(testCase.path, testCase.recursive, testCase.longFormat, testCase.showHidden, testCase.reverse, testCase.sortByModTime)
		
		// Stop capturing stdout
		w.Close()
		os.Stdout = old

		var buf bytes.Buffer
		io.Copy(&buf, r)
		actual := buf.String()

		if actual != testCase.expected {
			t.Errorf("For input %v, expected output %v, but got %v", testCase.path, testCase.expected, actual)
		}
	}
}
