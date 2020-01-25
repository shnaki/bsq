package main

import (
	"bytes"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func readMapString(file string) string {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	s, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return string(s)
}

func TestParseMap(t *testing.T) {
	testCases := []struct {
		name        string
		in          string
		errStr      string
		expectedMap Map
	}{
		{
			name: "example",
			in:   readMapString("./testdata/example_file"),
			expectedMap: Map{
				X:        27,
				Y:        9,
				Empty:    '.',
				Obstacle: 'o',
				Full:     'x',
				Obstacles: []*Obstacle{
					{Coordinate: Coordinate{X: 4, Y: 1}},
					{Coordinate: Coordinate{X: 12, Y: 2}},
					{Coordinate: Coordinate{X: 4, Y: 4}},
					{Coordinate: Coordinate{X: 15, Y: 5}},
					{Coordinate: Coordinate{X: 6, Y: 7}},
					{Coordinate: Coordinate{X: 21, Y: 7}},
					{Coordinate: Coordinate{X: 2, Y: 8}},
					{Coordinate: Coordinate{X: 10, Y: 8}},
				},
			},
		},
		{
			name:   "invalid header 1",
			in:     readMapString("./testdata/invalid-header-1"),
			errStr: "invalid header: 9.o",
		},
		{
			name:   "invalid header 2",
			in:     readMapString("./testdata/invalid-header-2"),
			errStr: `strconv.Atoi: parsing "........................": invalid syntax`,
		},
		{
			name:   "invalid line length",
			in:     readMapString("./testdata/invalid-line-length"),
			errStr: "line length is not 27 on line 4: .......................... ",
		},
		{
			name:   "contains full character",
			in:     readMapString("./testdata/full"),
			errStr: "full character is not allowed as input: xx.........................",
		},
		{
			name:   "invalid map character",
			in:     readMapString("./testdata/invalid-map-char"),
			errStr: "invalid map character a, candidates: [., o, x]",
		},
		{
			name:   "lacked data",
			in:     readMapString("./testdata/lacked-data"),
			errStr: "no map data",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			in := bytes.NewBufferString(tc.in)
			m, err := ParseMap(in)
			if tc.errStr != "" {
				if err == nil {
					t.Errorf("got %v, want %s", nil, tc.errStr)
				} else {
					errStr := err.Error()
					if errStr != tc.errStr {
						t.Errorf("got %v, want %s", errStr, tc.errStr)
					}
				}
			} else {
				if err != nil {
					t.Errorf("got %v, want %v", err, nil)
				}
			}
			if m == nil {
				return
			}
			if m.X != tc.expectedMap.X {
				t.Errorf("x: got %d, want %d", m.X, tc.expectedMap.X)
			}
			if m.Y != tc.expectedMap.Y {
				t.Errorf("y: got %d, want %d", m.Y, tc.expectedMap.Y)
			}
			if m.Empty != tc.expectedMap.Empty {
				t.Errorf("empty: got %s, want %s",
					string(m.Empty), string(tc.expectedMap.Empty))
			}
			if m.Obstacle != tc.expectedMap.Obstacle {
				t.Errorf("empty: got %s, want %s",
					string(m.Obstacle), string(tc.expectedMap.Obstacle))
			}
			if m.Full != tc.expectedMap.Full {
				t.Errorf("empty: got %s, want %s",
					string(m.Full), string(tc.expectedMap.Full))
			}
			lp := len(m.Obstacles)
			le := len(tc.expectedMap.Obstacles)
			if lp != le {
				t.Errorf("got %d obstacles, want %d obstacles", lp, le)
				return
			}
			for i, o := range m.Obstacles {
				eo := tc.expectedMap.Obstacles[i]
				if !reflect.DeepEqual(o, eo) {
					t.Errorf("got %+v, want %+v", o, eo)
				}
			}
		})
	}
}
