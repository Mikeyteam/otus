package main

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	var returnCode int
	var err error
	os.Setenv("HELLO", "SHOULD_REPLACE")
	os.Setenv("FOO", "SHOULD_REPLACE")
	os.Setenv("UNSET", "SHOULD_REMOVE")
	os.Setenv("ADDED", "from original env")
	os.Setenv("EMPTY", "SHOULD_BE_EMPTY")

	type args struct {
		cmd []string
		env Environment
	}
	tests := []struct {
		name           string
		args           args
		wantReturnCode int
		wantExpected   string
	}{
		{
			name: "original env",
			args: args{
				cmd: []string{"/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"},
				env: Environment{},
			},
			wantReturnCode: 0,
			wantExpected: `HELLO is (SHOULD_REPLACE)
BAR is ()
FOO is (SHOULD_REPLACE)
UNSET is (SHOULD_REMOVE)
ADDED is (from original env)
EMPTY is (SHOULD_BE_EMPTY)
arguments are arg1=1 arg2=2`,
		},
		{
			name: "echo.sh",
			args: args{
				cmd: []string{"/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2"},
				env: Environment{
					"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
					"BAR":   EnvValue{Value: "bar", NeedRemove: false},
					"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
					"EMPTY": EnvValue{Value: "", NeedRemove: false},
					"UNSET": EnvValue{Value: "", NeedRemove: true},
				},
			},
			wantReturnCode: 0,
			wantExpected: `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2`,
		},
		{
			name: "dont have cmd",
			args: args{
				cmd: []string{},
				env: Environment{
					"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
					"BAR":   EnvValue{Value: "bar", NeedRemove: false},
					"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
					"EMPTY": EnvValue{Value: "", NeedRemove: false},
					"UNSET": EnvValue{Value: "", NeedRemove: true},
				},
			},
			wantReturnCode: 1,
			wantExpected:   "",
		},
		{
			name: "empty env",
			args: args{
				cmd: []string{},
				env: Environment{},
			},
			wantReturnCode: 1,
			wantExpected:   "",
		},
		{
			name: "more 2 args",
			args: args{
				cmd: []string{"/bin/bash", "testdata/echo.sh", "arg1=1", "arg2=2", "arg3=3"},
				env: Environment{
					"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
					"BAR":   EnvValue{Value: "bar", NeedRemove: false},
					"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
					"EMPTY": EnvValue{Value: "", NeedRemove: false},
					"UNSET": EnvValue{Value: "", NeedRemove: true},
				},
			},
			wantReturnCode: 0,
			wantExpected: `HELLO is ("hello")
BAR is (bar)
FOO is (   foo
with new line)
UNSET is ()
ADDED is (from original env)
EMPTY is ()
arguments are arg1=1 arg2=2 arg3=3`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buffer bytes.Buffer
			var file *os.File
			func() {
				var fileInner *os.File
				file, fileInner, err = os.Pipe()
				defer fileInner.Close()
				origStdout := os.Stdout
				defer func() {
					os.Stdout = origStdout
				}()
				os.Stdout = fileInner
				require.NoError(t, err)
				returnCode = RunCmd(test.args.cmd, test.args.env)
			}()
			_, err = buffer.ReadFrom(file)
			require.NoError(t, err)
			require.Equalf(t, test.wantReturnCode, returnCode, "ReturnCode must be %v", test.wantReturnCode)
			require.Equal(t, test.wantExpected, strings.TrimRight(buffer.String(), "\n"), "must be equal")
		})
	}
}
