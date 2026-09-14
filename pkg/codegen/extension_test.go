package codegen

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_extTypeName(t *testing.T) {
	type args struct {
		extPropValue interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "success",
			args:    args{json.RawMessage(`"uint64"`)},
			want:    "uint64",
			wantErr: false,
		},
		{
			// kin-openapi hands us extension values already decoded.
			name:    "decoded value",
			args:    args{"uint64"},
			want:    "uint64",
			wantErr: false,
		},
		{
			name:    "wrong decoded type",
			args:    args{true},
			want:    "",
			wantErr: true,
		},
		{
			name:    "type conversion error",
			args:    args{nil},
			want:    "",
			wantErr: true,
		},
		{
			name:    "json unmarshal error",
			args:    args{json.RawMessage("invalid json format")},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extTypeName(tt.args.extPropValue)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_extParseBools(t *testing.T) {
	tests := []struct {
		name         string
		extPropValue interface{}
		want         bool
		wantErr      bool
	}{
		{name: "raw json", extPropValue: json.RawMessage(`true`), want: true},
		{name: "decoded value", extPropValue: true, want: true},
		{name: "decoded false", extPropValue: false, want: false},
		{name: "type conversion error", extPropValue: nil, wantErr: true},
		{name: "wrong decoded type", extPropValue: "true", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for name, parse := range map[string]func(interface{}) (bool, error){
				extPropOmitEmpty:                 extParseOmitEmpty,
				extPropGoTypeSkipOptionalPointer: extParsePropGoTypeSkipOptionalPointer,
			} {
				got, err := parse(tt.extPropValue)
				if tt.wantErr {
					assert.Errorf(t, err, "%s", name)
					continue
				}
				assert.NoErrorf(t, err, "%s", name)
				assert.Equalf(t, tt.want, got, "%s", name)
			}
		})
	}
}
