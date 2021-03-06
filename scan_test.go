/*-
 * Copyright (c) 2017, 2018
 *	Benny Siegert <bsiegert@gmail.com>
 *
 * Provided that these terms and disclaimer and all copyright notices
 * are retained or reproduced in an accompanying document, permission
 * is granted to deal in this work without restriction, including un-
 * limited rights to use, publicly perform, distribute, sell, modify,
 * merge, give away, or sublicence.
 *
 * This work is provided "AS IS" and WITHOUT WARRANTY of any kind, to
 * the utmost extent permitted by applicable law, neither express nor
 * implied; without malicious intent or gross negligence. In no event
 * may a licensor, author or contributor be held liable for indirect,
 * direct, other damage, loss, or other issues arising in any way out
 * of dealing in the work, even if advised of the possibility of such
 * damage or existence of a defect, except proven that it results out
 * of said person's immediate fault when using the work as intended.
 */

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
	{
		"#GO_SRCPATH=	github.com/foo/bar",
		"GO_SRCPATH",
		"",
	},
	{
		"foo\nGO_SRCPATH=	github.com/foo/bar\n\nbaz",
		"GO_SRCPATH",
		"github.com/foo/bar",
	},
	{
		"testing123",
		"GO_SRCPATH",
		"",
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
