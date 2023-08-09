package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestContainsInSlice(t *testing.T) {
	type args struct {
		items []string
		item  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "string in slice",
			args: args{
				items: []string{"a", "b", "c"},
				item:  "a",
			},
			want: true,
		},
		{
			name: "string not in slice",
			args: args{
				items: []string{"a", "b", "c"},
				item:  "x",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equalf(t, tt.want, ContainsInSlice(tt.args.items, tt.args.item), "ContainsInSlice(%v, %v)", tt.args.items, tt.args.item)
		})
	}
}
