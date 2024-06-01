package main

import (
	"fmt"
	"testing"
)

func Test_runDFA(t *testing.T) {
	type args struct {
		regex string
		str   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{

		{
			name: "Test 1",
			args: args{
				regex: "a",
				str:   "a",
			},
			want: true,
		},
		{
			name: "Test 2",
			args: args{
				regex: "a",
				str:   "b",
			},
			want: false,
		},
		{
			name: "Test 3",
			args: args{
				regex: "ab",
				str:   "ab",
			},
			want: true,
		},
		{
			name: "Test 4",
			args: args{
				regex: "ab",
				str:   "ba",
			},
			want: false,
		},
		{
			name: "Test 5",
			args: args{
				regex: "a*",
				str:   "aaaa",
			},
			want: true,
		},
		{
			name: "Test 6",
			args: args{
				regex: "a*",
				str:   "aaab",
			},
			want: false,
		},
		{
			name: "Test 7",
			args: args{
				regex: "(a|b)",
				str:   "a",
			},
			want: true,
		},
		{
			name: "Test 8",
			args: args{
				regex: "(a|b)",
				str:   "b",
			},
			want: true,
		},
		{
			name: "Test 9",
			args: args{
				regex: "(a|b)*",
				str:   "aaaaaaabbbbbb",
			},
			want: true,
		},
		{
			name: "Test 10",
			args: args{
				regex: "aaabbbb",
				str:   "aaabbbb",
			},
			want: true,
		},

		{
			name: "Test 11",
			args: args{
				regex: "aaba",
				str:   "aaba",
			},
			want: true,
		},

		{
			name: "Test 11",
			args: args{
				regex: "(a|b)*ba",
				str:   "aaaaaba",
			},
			want: true,
		},
		{
			name: "Test 11",
			args: args{
				regex: "ba(a|b)a",
				str:   "baba",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fmt.Printf("regex:%s str:%s\n", tt.args.regex, tt.args.str)
			if got := runDFA(tt.args.regex, tt.args.str); got != tt.want {
				t.Errorf("runDFA() = %v, want %v", got, tt.want)
			}
		})
	}
}
