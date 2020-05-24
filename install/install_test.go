package install

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func pstring(s string) *string {
	return &s
}

// IsFileExecutable determines the executable bits 0111 are set for the given file.
func IsFileExecutable(file string) (executable bool, err error) {
	stat, err := os.Stat(file)
	if err != nil {
		return
	}

	return stat.Mode().Perm()&0100 != 0, nil
}

func TestHooks(t *testing.T) {
	type args struct {
		preCommitContent  *string
		postCommitContent *string
		baseCommand       string
	}
	tests := []struct {
		name                  string
		args                  args
		wantPreCommitContent  string
		wantPostCommitContent string
	}{
		{
			name: "no hooks exist yet",
			args: args{
				preCommitContent:  nil,
				postCommitContent: nil,
				baseCommand:       "gogit",
			},
			wantPreCommitContent: `#!/bin/bash

gogit replace --replace-only-if-staged . # ` + defaultBashComment,
			wantPostCommitContent: `#!/bin/bash

gogit replace --replace-only-if-staged --undo . # ` + defaultBashComment,
		},

		{
			name: "add to existing pre-commit file with no match & empty file",
			args: args{
				preCommitContent:  pstring(""),
				postCommitContent: pstring(""),
				baseCommand:       "gogit",
			},
			wantPreCommitContent:  `gogit replace --replace-only-if-staged . # ` + defaultBashComment,
			wantPostCommitContent: `gogit replace --replace-only-if-staged --undo . # ` + defaultBashComment,
		},
		{
			name: "add to existing pre-commit file with no match & non-empty file",
			args: args{
				preCommitContent:  pstring("#!/bin/bash"),
				postCommitContent: pstring("#!/bin/bash"),
				baseCommand:       "gogit",
			},
			wantPreCommitContent: `#!/bin/bash

gogit replace --replace-only-if-staged . # ` + defaultBashComment,
			wantPostCommitContent: `#!/bin/bash

gogit replace --replace-only-if-staged --undo . # ` + defaultBashComment,
		},

		{
			name: "add to existing pre-commit file with match",
			args: args{
				preCommitContent:  pstring("# " + defaultBashComment),
				postCommitContent: pstring("# " + defaultBashComment),
				baseCommand:       "gogit",
			},
			wantPreCommitContent:  `gogit replace --replace-only-if-staged . # ` + defaultBashComment,
			wantPostCommitContent: `gogit replace --replace-only-if-staged --undo . # ` + defaultBashComment,
		},
	}

	for _, tt2 := range tests {
		t.Run(tt2.name, func(t *testing.T) {
			tt := tt2

			// Create the hooks folder.
			tempDir, err := ioutil.TempDir("", tt.name)
			if err != nil {
				t.Fatal(err)
			}

			t.Cleanup(func() {
				if err = os.RemoveAll(tempDir); err != nil {
					t.Fatal(err)
				}
			})

			base := tempDir
			if err = os.MkdirAll(filepath.Join(base, hooksPath()), 0755); err != nil {
				t.Fatal(err)
			}

			if tt.args.postCommitContent != nil {
				if err = ioutil.WriteFile(preCommitFilepath(base), []byte(*tt.args.preCommitContent), 0755); err != nil {
					t.Fatal(err)
				}
			}

			if tt.args.preCommitContent != nil {
				if err = ioutil.WriteFile(postCommitFilepath(base), []byte(*tt.args.postCommitContent), 0755); err != nil {
					t.Fatal(err)
				}
			}

			if err = Hooks(base, tt.args.baseCommand); err != nil {
				t.Fatal(err)
			}

			gotPreCommit, err := ioutil.ReadFile(preCommitFilepath(base))
			if err != nil {
				t.Fatal(err)
			}

			gotPostCommit, err := ioutil.ReadFile(postCommitFilepath(base))
			if err != nil {
				t.Fatal(err)
			}

			if !assert.Equal(t, tt.wantPreCommitContent, string(gotPreCommit), "pre-commit should be correct") {
				return
			}

			if !assert.Equal(t, tt.wantPostCommitContent, string(gotPostCommit), "pre-commit should be correct") {
				return
			}

			exec, err := IsFileExecutable(preCommitFilepath(base))
			if err != nil {
				t.Fatal(err)
			}

			if !exec {
				t.Fatal("pre-commit file should be executable")
			}

			exec, err = IsFileExecutable(postCommitFilepath(base))
			if err != nil {
				t.Fatal(err)
			}

			if !exec {
				t.Fatal("post-commit file should be executable")
			}
		})
	}
}
