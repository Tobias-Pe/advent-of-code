package main

import (
	"fmt"
	"testing"
)

func Test_isSafeWithRemovalOnce(t *testing.T) {
	tests := []struct {
		levels []int
		want   bool
	}{
		{
			levels: []int{7, 6, 4, 2, 1},
			want:   true,
		},
		{
			levels: []int{1, 2, 7, 8, 9},
			want:   false,
		},
		{
			levels: []int{1, 2, 3, 4, 9},
			want:   true,
		},
		{
			levels: []int{1, 2, 300, 4, 5},
			want:   true,
		},
		{
			levels: []int{1, 2, -300, 4, 5},
			want:   true,
		},
		{
			levels: []int{9, 7, 6, 2, 1},
			want:   false,
		},
		{
			levels: []int{1, 3, 2, 4, 5},
			want:   true,
		},
		{
			levels: []int{1, 3, 2, 1, 0},
			want:   true,
		},
		{
			levels: []int{8, 6, 4, 4, 1},
			want:   true,
		},
		{
			levels: []int{1, 3, 6, 7, 9},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.levels), func(t *testing.T) {
			if got := isSafeWithRemovalOnce(tt.levels); got != tt.want {
				t.Errorf("isSafeWithRemovalOnce() = %v, want %v: %v", got, tt.want, tt.levels)
			}
		})
	}
}
