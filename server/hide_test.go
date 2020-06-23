package server

import (
	"path/filepath"
	"testing"
)

func TestHide(t *testing.T) {

	hides := []string{
		// `.*`,
		// `*.git`,
		`/dpan/webdav/.git/*`,
	}

	testData := []string{
		`.git`,
		`a.git`,
		`b/a.git`,
		`/dpan/webdav/.git`,
		`/dpan/webdav/abc/.git`,
		`/dpan/webdav/.git`,
		`/dpan/webdav/.gita`,
		`/dpan/webdav/.gita/b`,
		`/dpan/webdav/.git/a`,
	}

	for _, hide := range hides {
		t.Log("-------", hide)
		for _, testPath := range testData {

			//matchOk, matchErr := path.Match(hide, testPath)

			matchOk, matchErr := filepath.Match(hide, testPath)
			t.Log(matchOk, matchErr)

		}
	}
}

func TestPwd(t *testing.T) {

	t.Log(UserPwd("a", "b", "d"))
}
