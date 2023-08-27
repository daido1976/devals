package main

import (
	"bytes"
	"testing"
)

func Test_run(t *testing.T) {
	type args struct {
		input        string
		keepComments bool
	}
	tests := []struct {
		name       string
		args       args
		wantOutput string
		wantErr    bool
	}{
		{
			name: "Test with keeping comments",
			args: args{
				input: `
# some comments
foo="ref+echo://foo"
bar="bar"
# baz="ref+echo://baz"

foobar=ref+echo://foo/bar
empty=""
none=

space="s p a c e"
`,
				keepComments: true,
			},
			wantOutput: `
# some comments
foo=foo
bar=bar
# baz="ref+echo://baz"

foobar=foo/bar
empty=
none=

space='s p a c e'
`,
			wantErr: false,
		},
		{
			name: "Test without keeping comments",
			args: args{
				input: `
# some comments
foo="ref+echo://foo"
`,
				keepComments: false,
			},
			wantOutput: `foo=foo
`,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			output := &bytes.Buffer{}
			if err := run(tt.args.input, output, tt.args.keepComments); (err != nil) != tt.wantErr {
				t.Errorf("run() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotOutput := output.String(); gotOutput != tt.wantOutput {
				t.Errorf("run() = %q, want %q", gotOutput, tt.wantOutput)
			}
		})
	}
}
