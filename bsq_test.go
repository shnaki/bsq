package main

import (
	"bytes"
	"testing"
)

func TestBsq(t *testing.T) {
	testCases := []struct {
		name        string
		in          string
		errStr      string
		expectedOut string
		expectedErr string
	}{
		{
			name: "example",
			in:   readMapString("./testdata/example_file"),
			expectedOut: `.....xxxxxxx...............
....oxxxxxxx...............
.....xxxxxxxo..............
.....xxxxxxx...............
....oxxxxxxx...............
.....xxxxxxx...o...........
.....xxxxxxx...............
......o..............o.....
..o.......o................
`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in := bytes.NewBufferString(tc.in)
			out := new(bytes.Buffer)
			errOut := new(bytes.Buffer)
			Bsq(in, out, errOut)
			str := out.String()
			if str != tc.expectedOut {
				t.Errorf("got %s, want %s", str, tc.expectedOut)
			}
		})
	}
}
