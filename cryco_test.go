package cryco

import (
	"io"
	"reflect"
	"testing"
)

func Test_exeName(t *testing.T) {
	tests := []struct {
		name    string
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := exeName()
			if (err != nil) != tt.wantErr {
				t.Errorf("exeName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("exeName() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetKey(t *testing.T) {
	tests := []struct {
		name    string
		want    []byte
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetKey()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDecrypt(t *testing.T) {
	type args struct {
		bKey      []byte
		cipherB64 string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decrypt(tt.args.bKey, tt.args.cipherB64)
			if (err != nil) != tt.wantErr {
				t.Errorf("Decrypt() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Decrypt() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_setValue(t *testing.T) {
	type args struct {
		p     interface{}
		field string
		value string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := setValue(tt.args.p, tt.args.field, tt.args.value); (err != nil) != tt.wantErr {
				t.Errorf("setValue() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetDefaults(t *testing.T) {
	type args struct {
		struc interface{}
		bKey  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := SetDefaults(tt.args.struc, tt.args.bKey); (err != nil) != tt.wantErr {
				t.Errorf("SetDefaults() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseReaders(t *testing.T) {
	type args struct {
		struc   interface{}
		readers []io.Reader
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseReaders(tt.args.struc, tt.args.readers); (err != nil) != tt.wantErr {
				t.Errorf("ParseReaders() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseFiles(t *testing.T) {
	type args struct {
		struc     interface{}
		filenames []string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ParseFiles(tt.args.struc, tt.args.filenames...); (err != nil) != tt.wantErr {
				t.Errorf("ParseFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
