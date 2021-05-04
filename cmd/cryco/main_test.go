package main

import (
	"testing"
)

func Test_allZero(t *testing.T) {
	type args struct {
		s []byte
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "empty", args: args{}, want: true},
		{name: "emptyx", args: args{[]byte{}}, want: true},
		{name: "1zero", args: args{[]byte{0}}, want: true},
		{name: "2zero", args: args{[]byte{0, 0}}, want: true},
		{name: "1oneEnd", args: args{[]byte{0, 1}}, want: false},
		{name: "2one", args: args{[]byte{1, 1}}, want: false},
		{name: "1oneStart", args: args{[]byte{1, 0}}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := allZero(tt.args.s); got != tt.want {
				t.Errorf("allZero() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_main(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			main()
		})
	}
}

func TestGenerateKey(t *testing.T) {
	type args struct {
		l int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateKey(tt.args.l); got != tt.want {
				t.Errorf("GenerateKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
