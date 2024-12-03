package main

import (
	"reflect"
	"testing"
)

func Test_findMulOperation(t *testing.T) {
	type args struct {
		input       string
		startOffset int
	}
	tests := []struct {
		name               string
		args               args
		wantFoundOperation *mulOperation
		wantSkipAmount     int
	}{
		{
			"Find first mul",
			args{
				input:       "xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))",
				startOffset: 1,
			},
			&mulOperation{
				left:  2,
				right: 4,
			},
			8,
		},
		{
			"No mul",
			args{
				input:       "xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))",
				startOffset: 0,
			},
			nil,
			1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFoundOperation, gotSkipAmount := findMulOperation(tt.args.input, tt.args.startOffset)
			if !reflect.DeepEqual(gotFoundOperation, tt.wantFoundOperation) {
				t.Errorf("findMulOperation() gotFoundOperation = %v, want %v", gotFoundOperation, tt.wantFoundOperation)
			}
			if gotSkipAmount != tt.wantSkipAmount {
				t.Errorf("findMulOperation() gotSkipAmount = %v, want %v", gotSkipAmount, tt.wantSkipAmount)
			}
		})
	}
}

func Test_sumAllMulOperations(t *testing.T) {
	type args struct {
		lines []string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"Find and sum muls, example input",
			args{
				lines: []string{"xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))"},
			},
			161,
		},
		{
			"Find and sum muls, example input with break",
			args{
				lines: []string{"xmul(2,4)%&mul[3,7]!@^do_not_mu", "l(5,5)+mul(32,64]then(mul(11,8)mul(8,5))"},
			},
			161,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := sumAllMulOperations(tt.args.lines); got != tt.want {
				t.Errorf("sumAllMulOperations() = %v, want %v", got, tt.want)
			}
		})
	}
}
