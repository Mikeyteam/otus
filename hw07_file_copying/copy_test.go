package main

import (
	"bytes"
	"errors"
	"os"
	"path"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopyNoError(t *testing.T) {
	tempDirectory, err := os.MkdirTemp("testdata", "tmp_")
	require.NoErrorf(t, err, "Couldn't create tempDirectory")
	defer os.RemoveAll(tempDirectory)

	type args struct {
		fromPath      string
		toPath        string
		offset        int64
		limit         int64
		compareSample string
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "offset0_limit0",
			args: args{
				fromPath:      path.Join("testdata", "input.txt"),
				toPath:        path.Join(tempDirectory, "result.txt"),
				offset:        0,
				limit:         0,
				compareSample: path.Join("testdata", "out_offset0_limit0.txt"),
			},
			wantErr: false,
		},
		{
			name: "offset0_limit10",
			args: args{
				fromPath:      path.Join("testdata", "input.txt"),
				toPath:        path.Join(tempDirectory, "result.txt"),
				offset:        0,
				limit:         10,
				compareSample: path.Join("testdata", "out_offset0_limit10.txt"),
			},
			wantErr: false,
		},
		{
			name: "offset0_limit1000",
			args: args{
				fromPath:      path.Join("testdata", "input.txt"),
				toPath:        path.Join(tempDirectory, "result.txt"),
				offset:        0,
				limit:         1000,
				compareSample: path.Join("testdata", "out_offset0_limit1000.txt"),
			},
			wantErr: false,
		},
		{
			name: "offset100_limit1000",
			args: args{
				fromPath:      path.Join("testdata", "input.txt"),
				toPath:        path.Join(tempDirectory, "result.txt"),
				offset:        100,
				limit:         1000,
				compareSample: path.Join("testdata", "out_offset100_limit1000.txt"),
			},
			wantErr: false,
		},
		{
			name: "offset6000_limit1000",
			args: args{
				fromPath:      path.Join("testdata", "input.txt"),
				toPath:        path.Join(tempDirectory, "result.txt"),
				offset:        6000,
				limit:         1000,
				compareSample: path.Join("testdata", "out_offset6000_limit1000.txt"),
			},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.False(t, test.wantErr, "Ok")
			compareSampleBytes, err := os.ReadFile(test.args.compareSample)
			require.NoErrorf(t, err, "Couldn't read file %v", test.args.compareSample)
			err = Copy(test.args.fromPath, test.args.toPath, test.args.offset, test.args.limit)
			require.NoErrorf(t, err, "Must be no error")
			compareResultBytes, err := os.ReadFile(test.args.toPath)
			require.NoError(t, err, "Dont read copied file %v", test.args.toPath)
			require.Truef(t,
				bytes.Equal(compareSampleBytes, compareResultBytes),
				"Content must be equal %v and %v",
				test.args.compareSample, test.args.toPath)
		})
	}
}

func TestCopyWithErrors(t *testing.T) {
	tempDirectory, err := os.MkdirTemp("testdata", "tmp_")
	require.NoErrorf(t, err, "Couldn't create tempDirectory")
	defer os.RemoveAll(tempDirectory)

	type args struct {
		fromPath  string
		toPath    string
		offset    int64
		limit     int64
		errorType error
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Error with offset",
			args: args{
				fromPath:  path.Join("testdata", "out_offset0_limit10.txt"),
				toPath:    path.Join(tempDirectory, "offset101_limit0"),
				offset:    101,
				limit:     0,
				errorType: ErrOffsetExceedsFileSize,
			},
			wantErr: true,
		},
		{
			name: "Error unsupported file",
			args: args{
				fromPath:  "/dev/urandom",
				toPath:    path.Join(tempDirectory, "offset0_limit100"),
				offset:    0,
				limit:     100,
				errorType: ErrUnsupportedFile,
			},
			wantErr: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			require.True(t, test.wantErr)
			err = Copy(test.args.fromPath, test.args.toPath, test.args.offset, test.args.limit)
			require.Error(t, err)
			require.Truef(t, errors.Is(err, test.args.errorType), "error string %q\nerror type %T", err, err)
		})
	}
}

func TestCopyProgressBar(t *testing.T) {
	tempDirectory, err := os.MkdirTemp("testdata", "tmp_")
	require.NoErrorf(t, err, "Couldn't create tempDirectory")
	defer os.RemoveAll(tempDirectory)

	type args struct {
		fromPath          string
		toPath            string
		offset            int64
		limit             int64
		progressBarSample string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "offset100_limit1000",
			args: args{
				fromPath:          path.Join("testdata", "input.txt"),
				toPath:            path.Join(tempDirectory, "result.txt"),
				offset:            100,
				limit:             1000,
				progressBarSample: "100.00%",
			},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buffer bytes.Buffer
			var file *os.File
			var err error
			func() {
				var w *os.File
				file, w, err = os.Pipe()
				defer w.Close()
				origStdout := os.Stdout
				defer func() { os.Stdout = origStdout }()
				os.Stdout = w
				require.NoError(t, err)
				err = Copy(test.args.fromPath, test.args.toPath, test.args.offset, test.args.limit)
			}()
			_, err = buffer.ReadFrom(file)
			stdOutResult := buffer.String()
			require.Truef(t,
				strings.Contains(stdOutResult, test.args.progressBarSample),
				"Result:\n %v \nBut should contains:\n%v", stdOutResult, test.args.progressBarSample)
		})
	}
}
