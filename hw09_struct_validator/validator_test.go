package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:" in:admin,stuff"`
		Phones []string `validate:" len :11"`
		meta   json.RawMessage
	}
	App struct {
		Version string `validate:"len:5"`
	}
	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}
	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		wantErr     bool
		expectedErr error
	}{
		{
			in: User{
				ID:     "1",
				Name:   "Andrey Larkin",
				Age:    60,
				Email:  "andrey-larkin.ru",
				Role:   "worker",
				Phones: []string{"012345543210", "454545"},
				meta:   []byte{},
			},
			wantErr: true,
			expectedErr: ErrorsValidation{
				ErrorValidation{
					Field: "ID",
					Err:   LenError{Limit: 36, CurrentValue: 1},
				},
				ErrorValidation{
					Field: "Age",
					Err:   MaxError{Limit: 50, CurrentValue: 60},
				},
				ErrorValidation{
					Field: "Email",
					Err:   ErrorRegexp,
				},
				ErrorValidation{
					Field: "Role",
					Err:   InError{Limit: "admin,stuff", CurrentValue: "worker"},
				},
				ErrorValidation{
					Field: "Phones",
					Err:   LenError{Limit: 11, CurrentValue: 12},
				},
			},
		},
		{
			in: User{
				ID:     "589347589347598347598375938753493444",
				Name:   "Andrey Larkin",
				Age:    50,
				Email:  "andrey@larkin.ru",
				Role:   "admin",
				Phones: []string{"23232323222", "33334444344"},
				meta:   []byte{},
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			in: User{
				ID:     "589347589347598347598375938753493444",
				Name:   "Andrey Larkin",
				Age:    10,
				Email:  "andrey@larkin.ru",
				Role:   "admin",
				Phones: []string{"01234567891", "12312312311"},
				meta:   []byte{},
			},
			wantErr: true,
			expectedErr: ErrorsValidation{
				ErrorValidation{
					Field: "Age",
					Err:   MinError{Limit: 18, CurrentValue: 10},
				},
			},
		},
		{
			in: App{
				Version: "v.1.2.3",
			},
			wantErr: true,
			expectedErr: ErrorsValidation{
				ErrorValidation{
					Field: "Version",
					Err:   LenError{Limit: 5, CurrentValue: 7},
				},
			},
		},
		{
			in: App{
				Version: "v.0.1",
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			in: Token{
				Header:    []byte("header"),
				Payload:   []byte("body"),
				Signature: []byte("parametr"),
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 200,
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			in: Response{
				Code: 301,
			},
			wantErr: true,
			expectedErr: ErrorsValidation{
				ErrorValidation{
					Field: "Code",
					Err:   InError{Limit: "200,404,500", CurrentValue: "301"},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			var valErr ErrorsValidation
			var expectedErr ErrorsValidation
			actualErr := Validate(tt.in)
			switch {
			case tt.wantErr:
				require.Error(t, actualErr)
				require.ErrorAs(t, actualErr, &valErr)
				require.ErrorAs(t, tt.expectedErr, &expectedErr)
				sort.Slice(valErr, func(i, j int) bool {
					if valErr[i].Field == valErr[j].Field {
						return valErr[i].Err.Error() < valErr[j].Err.Error()
					}
					return valErr[i].Field > valErr[j].Field
				})
				sort.Slice(expectedErr, func(i, j int) bool {
					if expectedErr[i].Field == expectedErr[j].Field {
						return expectedErr[i].Err.Error() < expectedErr[j].Err.Error()
					}
					return expectedErr[i].Field > expectedErr[j].Field
				})
				require.Equal(t, expectedErr, valErr)
			default:
				require.NoError(t, actualErr)
			}
		})
	}
}

func TestValidateSlices(t *testing.T) {
	type Slices struct {
		StringsCheckIn  []string `validate:"in:andrey1,andrey,andrey2"`
		StringsCheckLen []string `validate:"len:4"`
		Emails          []string `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Ages            []int    `validate:"min:10|max:16"`
		Limits          []int    `validate:"in:4,5,6"`
	}
	tests := []struct {
		in          interface{}
		wantErr     bool
		expectedErr error
	}{
		{
			in: Slices{
				StringsCheckIn:  []string{"andrey", "andrey1", "andrey2", "andrey1"},
				StringsCheckLen: []string{"1234", "4321", "abcd"},
				Emails:          []string{"mail@mail.ru", "mail1@mail.ru", "mail2@mail.ru"},
				Ages:            []int{11, 12, 13, 14},
				Limits:          []int{6, 6, 4, 4, 5},
			},
			wantErr:     false,
			expectedErr: nil,
		},
		{
			in: Slices{
				StringsCheckIn:  []string{"olol", "lol", "trololo", "andrey1"},
				StringsCheckLen: []string{"123", "432", "abc"},
				Emails:          []string{"mail-mail.ru", "mail1-mail.ru", "mail2-mail.ru"},
				Ages:            []int{8, 18},
				Limits:          []int{1, 2, 3},
			},
			wantErr: true,
			expectedErr: ErrorsValidation{
				ErrorValidation{
					Field: "StringsCheckIn", Err: InError{Limit: "andrey1,andrey,andrey2", CurrentValue: "olol"},
				},
				ErrorValidation{
					Field: "StringsCheckLen", Err: LenError{Limit: 4, CurrentValue: 3},
				},
				ErrorValidation{
					Field: "Emails", Err: ErrorRegexp,
				},
				ErrorValidation{
					Field: "Ages", Err: MinError{Limit: 10, CurrentValue: 8},
				},
				ErrorValidation{
					Field: "Ages", Err: MaxError{Limit: 16, CurrentValue: 18},
				},
				ErrorValidation{
					Field: "Limits", Err: InError{Limit: "4,5,6", CurrentValue: "1"},
				},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			var valErr ErrorsValidation
			var expectedErr ErrorsValidation
			actualErr := Validate(tt.in)
			switch {
			case tt.wantErr:
				require.Error(t, actualErr)
				require.ErrorAs(t, actualErr, &valErr)
				require.ErrorAs(t, tt.expectedErr, &expectedErr)
				sort.Slice(valErr, func(i, j int) bool {
					if valErr[i].Field == valErr[j].Field {
						return valErr[i].Err.Error() < valErr[j].Err.Error()
					}
					return valErr[i].Field > valErr[j].Field
				})
				sort.Slice(expectedErr, func(i, j int) bool {
					if expectedErr[i].Field == expectedErr[j].Field {
						return expectedErr[i].Err.Error() < expectedErr[j].Err.Error()
					}
					return expectedErr[i].Field > expectedErr[j].Field
				})
				require.Equal(t, expectedErr, valErr)
			default:
				require.NoError(t, actualErr)
			}
		})
	}
}
