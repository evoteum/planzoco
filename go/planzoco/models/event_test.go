package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWinningOptions(t *testing.T) {
	tests := []struct {
		name     string
		question Question
		want     []Option
	}{
		{
			name: "no options",
			question: Question{
				Options: []Option{},
			},
			want: nil,
		},
		{
			name: "no votes",
			question: Question{
				Options: []Option{
					{Text: "Option 1", Votes: 0},
					{Text: "Option 2", Votes: 0},
				},
			},
			want: nil,
		},
		{
			name: "single winner",
			question: Question{
				Options: []Option{
					{Text: "Option 1", Votes: 1},
					{Text: "Option 2", Votes: 2},
					{Text: "Option 3", Votes: 0},
				},
			},
			want: []Option{
				{Text: "Option 2", Votes: 2},
			},
		},
		{
			name: "tied winners",
			question: Question{
				Options: []Option{
					{Text: "Option 1", Votes: 2},
					{Text: "Option 2", Votes: 2},
					{Text: "Option 3", Votes: 1},
				},
			},
			want: []Option{
				{Text: "Option 1", Votes: 2},
				{Text: "Option 2", Votes: 2},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.question.WinningOptions()
			assert.Equal(t, tt.want, got)
		})
	}
}
