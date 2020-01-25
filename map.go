package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

// Map はマップ情報を格納する構造体です。
type Map struct {
	X         int
	Y         int
	Empty     byte
	Obstacle  byte
	Full      byte
	Obstacles []*Obstacle
}

// Obstacle は障害物を表す構造体です。
type Obstacle struct {
	Coordinate
}

// Coordinate は(x, y)座標値を格納する構造体です。
type Coordinate struct {
	X int
	Y int
}

// NewMap はMapを初期化して返します。
func NewMap(y int, empty, obstacle, full byte) (*Map, error) {
	return &Map{
		Y:        y,
		Empty:    empty,
		Obstacle: obstacle,
		Full:     full,
	}, nil
}

// NewObstacle はObstacleを初期化して返します。
func NewObstacle(x, y int) (*Obstacle, error) {
	p := &Obstacle{
		Coordinate: Coordinate{
			X: x,
			Y: y,
		},
	}
	return p, nil
}

// ParseMap はReaderからMap情報を読み込んでMapを返します。
func ParseMap(r io.Reader) (*Map, error) {
	scanner := bufio.NewScanner(r)
	var i int
	var m *Map
	for scanner.Scan() {
		line := scanner.Text()
		if i == 0 {
			// 1行目は"9.ox"のようなヘッダー行。
			// フォーマット: (行数)(empty文字)(obstacle文字)(full文字)。
			var bytes = []byte(line)
			l := len(bytes)
			if l < 4 {
				// ヘッダーが不正の場合はmap errorを返す。
				return nil, fmt.Errorf("invalid header: %s", line)
			}
			s := string(bytes[0 : l-3])
			y, err := strconv.Atoi(s)
			if err != nil {
				return nil, err
			}
			full := bytes[l-1]
			obstacle := bytes[l-2]
			empty := bytes[l-3]
			m, err = NewMap(y, empty, obstacle, full)
			if err != nil {
				return nil, err
			}
		} else {
			// 2行目以降はマップ文字。
			if m == nil {
				return nil, errors.New("no map data")
			}
			ll := len(line)
			if i == 1 {
				m.X = ll
			} else {
				// マップの文字数が異なる場合はエラーを返す。
				if ll != m.X {
					return nil, fmt.Errorf("line length is not %d on line %d: %s ",
						m.X, i, line)
				}
			}

			// 障害物の位置をMapに格納する。
			bytes := []byte(line)
			for j, b := range bytes {
				switch b {
				case m.Empty:
				case m.Obstacle:
					x := j
					y := i - 1
					o, err := NewObstacle(x, y)
					if err != nil {
						return nil, err
					}
					m.Obstacles = append(m.Obstacles, o)
				case m.Full:
					return nil, fmt.Errorf("full character is not allowed as input: %s", line)
				default:
					return nil, fmt.Errorf("invalid map character %s, candidates: [%s, %s, %s]",
						string(b), string(m.Empty), string(m.Obstacle), string(m.Full))
				}
			}
		}
		if m != nil && i == m.Y {
			return m, nil
		}
		i++
	}
	return nil, errors.New("no map data")
}
