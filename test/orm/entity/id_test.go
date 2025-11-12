package entity_test

import (
	"fmt"
	"testing"

	. "github.com/philiphil/restman/orm/entity"
)

func TestId_String(t *testing.T) {
	tests := []struct {
		name string
		id   ID
		want string
	}{
		{
			name: "1",
			id:   1,
			want: "1",
		},
		{
			name: "2",
			id:   2,
			want: "2",
		},
		{
			name: "3",
			id:   3,
			want: "3",
		},
		{
			name: "4",
			id:   4,
			want: "4",
		},
		{
			name: "5",
			id:   5,
			want: "5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CastId(tt.id); fmt.Sprintf("%v", got) != tt.want {
				t.Errorf("CastId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCastId(t *testing.T) {
	tests := []struct {
		name string
		id   any
		want ID
	}{
		{
			name: "1",
			id:   1,
			want: 1,
		},
		{
			name: "2",
			id:   uint(2),
			want: 2,
		},
		{
			name: "3",
			id:   "3",
			want: 3,
		},
		{
			name: "4",
			id:   ID(4),
			want: 4,
		},
		{
			name: "5",
			id:   error(nil),
			want: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CastId(tt.id); got != tt.want {
				t.Errorf("CastId() = %v, want %v", got, tt.want)
			}
		})
	}
}
