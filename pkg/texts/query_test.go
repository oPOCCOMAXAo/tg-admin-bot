package texts

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDecodeQuery(t *testing.T) {
	testCases := []struct {
		input  string
		output Query
	}{
		{
			input: "",
			output: Query{
				Params: map[string][]string{},
			},
		},
		{
			input: "/add",
			output: Query{
				Command: "/add",
				Params: map[string][]string{
					"/add": nil,
				},
			},
		},
		{
			input: "multi word=2,3 params=4",
			output: Query{
				Command: "multi",
				Params: map[string][]string{
					"multi":  nil,
					"word":   {"2", "3"},
					"params": {"4"},
				},
			},
		},
		{
			input: "multi word=2 word=3 params",
			output: Query{
				Command: "multi",
				Params: map[string][]string{
					"multi":  nil,
					"word":   {"2", "3"},
					"params": nil,
				},
			},
		},
		{
			input: "command=test with params",
			output: Query{
				Command: "command",
				Params: map[string][]string{
					"command": {"test"},
					"with":    nil,
					"params":  nil,
				},
			},
		},
	}

	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			res := DecodeQuery(tC.input)
			require.Equal(t, tC.output, res)
		})
	}
}

func TestQuery_Encode(t *testing.T) {
	testCases := []struct {
		input  Query
		output string
	}{
		{
			input: Query{
				Params: map[string][]string{},
			},
			output: "",
		},
		{
			input: Query{
				Params: map[string][]string{
					"add": nil,
				},
			},
			output: "add",
		},
		{
			input: Query{
				Command: "multi",
				Params: map[string][]string{
					"multi":  nil,
					"word":   {"2", "3"},
					"params": {"4"},
				},
			},
			output: "multi params=4 word=2,3",
		},
		{
			input: Query{
				Command: "command",
				Params:  map[string][]string{},
			},
			output: "command",
		},
		{
			input: Query{
				Command: "command",
				Params: map[string][]string{
					"command": {"test"},
					"with":    nil,
					"params":  nil,
				},
			},
			output: "command=test params with",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.output, func(t *testing.T) {
			res := tC.input.Encode()
			require.Equal(t, tC.output, res)
		})
	}
}
