package main

import "testing"

var extractVarTests = []struct {
	contents, varname, result string
}{
	{
		"GO_SRCPATH=	github.com/foo/bar",
		"GO_SRCPATH",
		"github.com/foo/bar",
	},
}

func TestExtractVar(t *testing.T) {
	for _, test := range extractVarTests {
		got, want := extractVar([]byte(test.contents), []byte(test.varname)), test.result
		if got != want {
			t.Errorf("extractVar(%q, %q): got %q, want %q", test.contents, test.varname, got, want)
		}
	}
}
