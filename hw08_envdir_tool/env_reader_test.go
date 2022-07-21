package main

import (
	"reflect"
	"testing"
)

func TestReadDir(t *testing.T) {
	type args struct {
		directory string
	}
	tests := []struct {
		name    string
		args    args
		wantSee Environment
		wantErr bool
	}{
		{
			name: "success directory",
			args: args{directory: "testdata/env"},
			wantSee: Environment{
				"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
				"BAR":   EnvValue{Value: "bar", NeedRemove: false},
				"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
				"EMPTY": EnvValue{Value: "", NeedRemove: false},
				"UNSET": EnvValue{Value: "", NeedRemove: true},
			},
			wantErr: false,
		},
		{
			name:    "error directory",
			args:    args{directory: "testdata/unknow"},
			wantSee: nil,
			wantErr: true,
		},
		{
			name: "empty directory",
			args: args{directory: "testdata/env"},
			wantSee: Environment{
				"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
				"BAR":   EnvValue{Value: "bar", NeedRemove: false},
				"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
				"EMPTY": EnvValue{Value: "", NeedRemove: false},
				"UNSET": EnvValue{Value: "", NeedRemove: true},
			},
			wantErr: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := ReadDir(test.args.directory)
			if (err != nil) != test.wantErr {
				t.Errorf("ReadDir() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.wantSee) {
				t.Errorf("ReadDir() = %v, wantSee %v", got, test.wantSee)
			}
		})
	}
}
