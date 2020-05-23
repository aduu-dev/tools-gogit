package gogit

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

type testdata struct {
	name  string
	input []byte
	want1 []byte
	want2 []byte
}

func readTestData(t *testing.T, base string) *testdata {
	input, err := ioutil.ReadFile(filepath.Join(base, "input"))
	assert.NoError(t, err)

	want1, err := ioutil.ReadFile(filepath.Join(base, "want1"))
	assert.NoError(t, err)

	want2, err := ioutil.ReadFile(filepath.Join(base, "want2"))
	assert.NoError(t, err)

	return &testdata{
		name:  filepath.Base(base),
		input: input,
		want1: want1,
		want2: want2,
	}
}

func (test *testdata) runTest(t *testing.T) {
	tempDir, err := ioutil.TempDir("", "e2e-test")
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err = os.RemoveAll(tempDir); err != nil {
			t.Fatal(err)
		}
	})

	path := tempDir
	goModFilepath := filepath.Join(tempDir, "go.mod")
	backup := filepath.Join(tempDir, backupFilename())
	if err = ioutil.WriteFile(goModFilepath, test.input, 0755); err != nil {
		t.Fatal(err)
	}

	if err = RemoveLocalReplacesFromGomod(path); err != nil {
		t.Fatal(err)
	}

	if !assert.FileExists(t, backup, "backup file was not created") {
		return
	}

	// Check the correct go.mod was written.
	got1, err := ioutil.ReadFile(goModFilepath)
	if err != nil {
		t.Fatal(err)
	}

	if !assert.Equal(t, string(test.want1), string(got1), "failed to remove local replace statements") {
		return
	}

	if err = UndoRemovingLocalReplacesFromGomod(path); err != nil {
		t.Fatal(err)
	}

	// Check the correct go.mod was written.
	got2, err := ioutil.ReadFile(goModFilepath)
	if err != nil {
		t.Fatal(err)
	}

	if !assert.Equal(t, string(test.want2), string(got2), "failed to undo the local changes") {
		return
	}

	if !assert.NoFileExists(t, backup, "backup was not removed") {
		return
	}
}

func TestEndToEnd(t *testing.T) {
	e2eFilepath := filepath.Join("testdata", "e2e")

	// Read test folders.
	dir, err := ioutil.ReadDir(e2eFilepath)
	if err != nil {
		return
	}

	// read tests.
	tests := make([]*testdata, 0, len(dir))
	for _, test := range dir {
		tests = append(tests, readTestData(t, filepath.Join(e2eFilepath, test.Name())))
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tt := test
			tt.runTest(t)
		})
	}
}

func TestRemoveLocalReplacesFromGomod_error_if_backupfile_exists(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), t.Name())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err = os.RemoveAll(tempDir); err != nil {
			t.Fatal(err)
		}
	})

	backup := filepath.Join(tempDir, backupFilename())

	// Write the backup file.
	if err = ioutil.WriteFile(backup, []byte{}, 0755); err != nil {
		t.Fatal(err)
	}

	// Write a dummy go.mod file.
	dummyGomodFile := `module aduu.dev/k

go 1.14

require (
	aduu.dev/utils v0.0.0-20200523102358-b59e2ccc9c3e
	fyne.io/fyne v1.2.4
	github.com/golang/protobuf v1.4.2
	github.com/google/go-github/v31 v31.0.0
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/spf13/cobra v1.0.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.23.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
	k8s.io/klog/v2 v2.0.0
)

replace aduu.dev/utils => ../aduu-dev-utils`

	path := tempDir
	goModFilepath := filepath.Join(tempDir, "go.mod")
	if err = ioutil.WriteFile(goModFilepath, []byte(dummyGomodFile), 0755); err != nil {
		t.Fatal(err)
	}

	if err = RemoveLocalReplacesFromGomod(path); err == nil {
		t.Fatalf("RemoveLocalReplaces should have returned error when there is no go.mod.b file")
	}
}

func Test(t *testing.T) {
	tempDir, err := ioutil.TempDir(os.TempDir(), t.Name())
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		if err = os.RemoveAll(tempDir); err != nil {
			t.Fatal(err)
		}
	})

	backup := filepath.Join(tempDir, backupFilename())

	// Write a dummy go.mod file.
	dummyGomodFile := `module aduu.dev/k

go 1.14

require (
	aduu.dev/utils v0.0.0-20200523102358-b59e2ccc9c3e
	fyne.io/fyne v1.2.4
	github.com/golang/protobuf v1.4.2
	github.com/google/go-github/v31 v31.0.0
	github.com/imdario/mergo v0.3.9 // indirect
	github.com/skratchdot/open-golang v0.0.0-20200116055534-eef842397966
	github.com/spf13/cobra v1.0.0
	golang.org/x/oauth2 v0.0.0-20200107190931-bf48bf16ab8d
	google.golang.org/grpc v1.29.1
	google.golang.org/protobuf v1.23.0
	k8s.io/apimachinery v0.17.0
	k8s.io/client-go v0.17.0
	k8s.io/klog/v2 v2.0.0
)

replace aduu.dev/utils => ../aduu-dev-utils`

	path := tempDir
	goModFilepath := filepath.Join(tempDir, "go.mod")
	if err = ioutil.WriteFile(goModFilepath, []byte(dummyGomodFile), 0755); err != nil {
		t.Fatal(err)
	}

	if err = RemoveLocalReplacesFromGomod(path); err != nil {
		t.Fatal(err)
	}

	// Remove backup.
	if err = os.Remove(backup); err != nil {
		t.Fatal(err)
	}

	if err = UndoRemovingLocalReplacesFromGomod(path); err == nil {
		t.Fatalf("Undo should return error if there is no go.mod.b backup file")
	}
}

func Test_removeLocalReplaceDirectives(t *testing.T) {
	type args struct {
		directives []*modfile.Replace
	}
	tests := []struct {
		name    string
		args    args
		wantOut []*modfile.Replace
	}{
		{
			name: "remove 1 local directive",
			args: args{
				directives: []*modfile.Replace{
					{
						Old: module.Version{
							Path:    "aduu.dev/utils",
							Version: "",
						},
						New: module.Version{
							Path:    "../aduu-dev-utils",
							Version: "",
						},
						Syntax: nil,
					},
				},
			},
			wantOut: []*modfile.Replace{
				{
					Old: module.Version{
						Path:    "aduu.dev/utils",
						Version: "",
					},
					New: module.Version{
						Path:    "../aduu-dev-utils",
						Version: "",
					},
					Syntax: nil,
				},
			},
		},
		{
			name: "do not remove non-local replace directives",
			args: args{
				directives: []*modfile.Replace{
					{
						Old: module.Version{
							Path:    "aduu.dev/utils",
							Version: "",
						},
						New: module.Version{
							Path:    "aduu.dev/hello/world",
							Version: "",
						},
						Syntax: nil,
					},
				},
			},
			wantOut: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotOut := removeLocalReplaceDirectives(tt.args.directives); !reflect.DeepEqual(gotOut, tt.wantOut) {
				t.Errorf("removeLocalReplaceDirectives() =\n got: %#v, \nwant: %#v", turnToNonPointerSlice(gotOut), turnToNonPointerSlice(tt.wantOut))
			}
		})
	}
}

func Test_isLocalDirective(t *testing.T) {
	type args struct {
		rep *modfile.Replace
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "is non-local directive",
			args: args{
				rep: &modfile.Replace{
					Old: module.Version{
						Path:    "aduu.dev/utils",
						Version: "",
					},
					New: module.Version{
						Path:    "aduu.dev/hello/world",
						Version: "",
					},
					Syntax: nil,
				},
			},
			want: false,
		},
		{
			name: "has ../ prefix",
			args: args{
				rep: &modfile.Replace{
					Old: module.Version{
						Path:    "aduu.dev/utils",
						Version: "",
					},
					New: module.Version{
						Path:    "../aduu-dev-utils",
						Version: "",
					},
					Syntax: nil,
				},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isLocalDirective(tt.args.rep); got != tt.want {
				t.Errorf("isLocalDirective(%#v) = %v, want %v", *tt.args.rep, got, tt.want)
			}
		})
	}
}
