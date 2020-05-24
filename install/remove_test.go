package install

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemove(t *testing.T) {
	type args struct {
		preCommit string
		postCommit string
	}

	tests := []struct {
		name string
		args args
		wantPreCommit string
		wantPostCommit string
	}{
		{
			name: "no match in either",
			args: args{
				preCommit:  `#!/bin/bash`,
				postCommit: ``,
			},
			wantPreCommit:  `#!/bin/bash`,
			wantPostCommit: ``,
		},
		{
			name: "match: remove empty line",
			args: args{
				preCommit:  `#!/bin/bash
# ` + defaultBashComment,
				postCommit: `# ` + defaultBashComment,
			},
			wantPreCommit:  `#!/bin/bash`,
			wantPostCommit: ``,
		},
		{
			name: "match: remove full line",
			args: args{
				preCommit:  `#!/bin/bash
 hello # ` + defaultBashComment,
				postCommit: ` hi# ` + defaultBashComment,
			},
			wantPreCommit:  `#!/bin/bash`,
			wantPostCommit: ``,
		},

		{
			name: "match: remove from the end of the file",
			args: args{
				preCommit:  `#!/bin/bash



 hello # ` + defaultBashComment,
				postCommit: ` hi# ` + defaultBashComment,
			},
			wantPreCommit:  `#!/bin/bash`,
			wantPostCommit: ``,
		},
	}

	for _, tt2 := range tests {
		t.Run(tt2.name, func(t *testing.T) {
			tt := tt2

			tempDir, err := ioutil.TempDir(os.TempDir(), strings.ReplaceAll(t.Name(), "/", "-"))
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() {
				if err = os.RemoveAll(tempDir); err != nil {
					t.Fatal(err)
				}
			})

			// Create hooks folder and write hooks content.
			base := tempDir
			if err = os.MkdirAll(filepath.Join(base, hooksPath()), 0755); err != nil {
				t.Fatal(err)
			}

			if err = ioutil.WriteFile(preCommitFilepath(base), []byte(tt.args.preCommit), 0755); err != nil {
				t.Fatal(err)
			}

			if err = ioutil.WriteFile(postCommitFilepath(base), []byte(tt.args.postCommit), 0755); err != nil {
				t.Fatal(err)
			}

			if err = Remove(base); err != nil {
				t.Fatal(err)
			}

			fileHasContent(t, preCommitFilepath(base), tt.wantPreCommit, "pre-commit should have this content")
			fileHasContent(t, postCommitFilepath(base), tt.wantPostCommit, "post-commit should have this content")
		})
	}
}

func fileHasContent(t *testing.T, file string, want string, msg string) {
	got, err := ioutil.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}

	if len(msg) != 0 {
		msg = "msg=" + msg
	}

	assert.Equalf(t, want, string(got), "file %#v should have the given content" + msg, file)
}

func TestRemove_returns_error_if_no_hook_files(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), strings.ReplaceAll(t.Name(), "/", "-"))
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err = os.RemoveAll(tempDir); err != nil {
			t.Fatal(err)
		}
	})

	// Create pre-commit and post-commit files.
	base := tempDir
	if err = os.MkdirAll(filepath.Join(base, hooksPath()), 0755); err != nil {
		t.Fatal(err)
	}

	if err = ioutil.WriteFile(preCommitFilepath(base), []byte(""), 0755); err != nil {
		t.Fatal(err)
	}

	if err = ioutil.WriteFile(postCommitFilepath(base), []byte(""), 0755); err != nil {
		t.Fatal(err)
	}

	if err = Remove(base); err != nil {
		t.Fatalf("error: Remove should not fail on empty commit hooks: %v", err)
	}

	// Check that if there is no match it does still not error.
	fileHasContent(t, preCommitFilepath(base), "", "")
	fileHasContent(t, postCommitFilepath(base), "", "")

	// Remove one file and see what happens.
	if err = os.Remove(preCommitFilepath(base)); err != nil {
		t.Fatal(err)
	}

	if err = Remove(base); err == nil {
		t.Fatalf("Remove should choke on non-existent pre-hook file.")
	}

	// Re-create pre-hook file.
	if err = ioutil.WriteFile(preCommitFilepath(base), []byte(""), 0755); err != nil {
		t.Fatal(err)
	}

	// Check pre-hook file got created correctly..
	fileHasContent(t, preCommitFilepath(base), "", "")

	// Remove post-commit file.
	if err = os.Remove(postCommitFilepath(base)); err != nil {
		t.Fatal(err)
	}

	if err = Remove(base); err == nil {
		t.Fatalf("Remove should choke on non-existent post-hook file.")
	}
}