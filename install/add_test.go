package install

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_addToShellFile_success(t *testing.T) {
	type args struct {
		content    string
		line    string
		comment string
	}
	successfulTests := []struct {
		name    string
		args    args
		want string
	}{
		// line is empty
		{
			name:    "line empty: existing matching => Replace with new comment",
			args:    args{
				content:    `
# AUTO-GENERATED
`,
				line:    "",
				comment: "AUTO-GEN",
			},
			want: `
# AUTO-GEN
`,
		},
		{
			name:    "line empty: existing matching exactly => Still remove spaces",
			args:    args{
				content:    `
  # AUTO-GEN
`,
				line:    "",
				comment: "AUTO-GEN",
			},
			want: `
# AUTO-GEN
`,
		},
		{
			name:    "line empty: no match => Add comment at the end with one empty line between prior content",
			args:    args{
				content:    `hello`,
				line:    "",
				comment: "AUTO-GEN",
			},
			want: `hello

# AUTO-GEN`,
		},

		// Line is set
		// line is empty
		{
			name:    "line is set: existing matching => Replace with new comment",
			args:    args{
				content:    `
# AUTO-GEN
`,
				line:    "my cmd",
				comment: "AUTO-GEN",
			},
			want: `
my cmd # AUTO-GEN
`,
		},
		{
			name:    "line is set: existing matching exactly => Still remove spaces",
			args:    args{
				content:    `
  # AUTO-GEN
`,
				line:    "my cmd",
				comment: "AUTO-GEN",
			},
			want: `
my cmd # AUTO-GEN
`,
		},
		{
			name:    "line is set: no match => Add comment at the end with one empty line between prior content",
			args:    args{
				content:    `hello`,
				line:    "my cmd",
				comment: "AUTO-GEN",
			},
			want: `hello

my cmd # AUTO-GEN`,
		},
	}

	for _, tt := range successfulTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := addToShellFile(tt.args.content, tt.args.line, tt.args.comment)
			if err != nil {
				t.Fatal(err)
			}

			if !assert.Equal(t, tt.want, got) {
				return
			}
		})
	}
}

func Test_addToShellFile_errors(t *testing.T) {
	type args struct {
		content    string
		line    string
		comment string
	}
	errorTests := []struct {
		name    string
		args    args
		wantErr error
	}{
		{
			name:    "line contains line separators",
			args:    args{
				content:    `
# AUTO-GENERATED by gogit`,
				line:    "\n",
				comment: "AUTO-GENERATED by gogit",
			},
			wantErr: errLineContainsLineSeparators,
		},
		{
			name:    "comment contains line separators",
			args:    args{
				content:    `
# AUTO-GENERATED by gogit`,
				line:    "",
				comment: "\nAUTO-GENERATED by gogit",
			},
			wantErr: errCommentContainsLineSeparators,
		},
		{
			name:    "comment is empty",
			args:    args{
				content:    ``,
				line:    "",
				comment: "",
			},
			wantErr: errCommentIsEmpty,
		},
		{
			name:    "error: comment contains hash",
			args:    args{
				content:    ``,
				line:    "",
				comment: "# abc",
			},
			wantErr: errCommentContainsHash,
		},
	}

	for _, tt := range errorTests {
		t.Run(tt.name, func(t *testing.T) {
			_, gotErr := addToShellFile(tt.args.content, tt.args.line, tt.args.comment)

			if !assert.EqualError(t, gotErr, tt.wantErr.Error()) {
				t.Errorf("addToShellFiles(%#v) =>\n gotErr: %v\nwantErr: %v\n", tt.args, gotErr, tt.wantErr)
			}
		})
	}
}

func Test_removeFromShellFile(t *testing.T) {
	type args struct {
		content string
		comment string
	}
	tests := []struct {
		name               string
		args               args
		wantChangedContent string
		wantErr            bool
	}{
		{
			name:               "two lines match the comment",
			args:               args{
				content: `
# AUTO-GEN
# AUTO-GENERATED`,
				comment: "AUTO-GEN",
			},
			wantChangedContent: ``,
			wantErr:            true,
		},
		{
			name:               "remove the single matching line",
			args:               args{
				content: `abc
# AUTO-GEN`,
				comment: "AUTO-GEN",
			},
			wantChangedContent: `abc`,
			wantErr:            false,
		},
		{
			name:               "no match: remove nothing",
			args:               args{
				content: `abc
# AUTO-GEN`,
				comment: "AUTO-GENE",
			},
			wantChangedContent: `abc
# AUTO-GEN`,
			wantErr:            false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotChangedContent, err := removeFromShellFile(tt.args.content, tt.args.comment)
			if (err != nil) != tt.wantErr {
				t.Errorf("removeFromShellFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotChangedContent != tt.wantChangedContent {
				t.Errorf("removeFromShellFile() gotChangedContent = %v, want %v", gotChangedContent, tt.wantChangedContent)
			}
		})
	}
}