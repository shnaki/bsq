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
	Square    *Square
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
func NewMap(y int, empty, obstacle, full byte) *Map {
	return &Map{
		Y:        y,
		Empty:    empty,
		Obstacle: obstacle,
		Full:     full,
	}
}

// NewObstacle はObstacleを初期化して返します。
func NewObstacle(x, y int) *Obstacle {
	return &Obstacle{
		Coordinate: Coordinate{
			X: x,
			Y: y,
		},
	}
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
			m = NewMap(y, empty, obstacle, full)
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
					m.AddObstacle(j, i-1)
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

// AddObstacle はマップに障害物を追加します。
func (m *Map) AddObstacle(x, y int) {
	o := NewObstacle(x, y)
	m.Obstacles = append(m.Obstacles, o)
}

// ObstacleGrid はbool型のマップサイズと同じサイズの二次元スライスを返します。
// 配列の障害物の位置はtrue, その他の位置はすべてfalseです。
func (m *Map) ObstacleGrid() [][]bool {
	g := make([][]bool, m.Y, m.Y)
	for i := range g {
		g[i] = make([]bool, m.X, m.X)
	}
	for _, o := range m.Obstacles {
		g[o.Y][o.X] = true
	}
	return g
}

// String はマップを文字列化して返します。
func (m *Map) String() string {
	var s string
	g := m.ObstacleGrid()
	for i, row := range g {
		for j, isObstacle := range row {
			if isObstacle {
				s += string(m.Obstacle)
				continue
			}

			// 正方形が設定されている場合はfull文字を出力する
			if m.Square != nil {
				if i >= m.Square.Start.Y &&
					j >= m.Square.Start.X &&
					i < m.Square.Start.Y+m.Square.Size &&
					j < m.Square.Start.X+m.Square.Size {
					s += string(m.Full)
					continue
				}
			}
			s += string(m.Empty)
		}
		s += "\n"
	}
	return s
}
