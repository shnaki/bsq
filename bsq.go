package main

import (
	"fmt"
	"io"
)

// Square は正方形の情報を格納する構造体です。
type Square struct {
	Start Coordinate
	Size  int
}

// NewSquare はSquareを初期化して返します。
func NewSquare(x, y, size int) *Square {
	return &Square{
		Start: Coordinate{
			X: x,
			Y: y,
		},
		Size: size,
	}
}

// Bsq はマップ情報を読み込んでoutに正方形を付加したマップ、errOutにエラー情報を出力します。
func Bsq(in io.Reader, out, errOut io.Writer) {
	m, err := ParseMap(in)
	if err != nil {
		fmt.Fprintln(errOut, "map error")
		return
	}
	g := m.ObstacleGrid()

	// グリッドの各点を正方形の左上の点とした場合のサイズを求める。
	var maxSize int
	for i, row := range g {
		for j, isObstacle := range row {
			// 始点が障害物ならサイズは0となるので無視する。
			if !isObstacle {
				s := NewSquare(j, i, 1)
				for isExpandable(s, g) {
					s.Size++
				}
				if s.Size > maxSize {
					maxSize = s.Size
					m.Square = s
				}
			}
		}
	}
	fmt.Fprintf(out, "%v", m)
}

func isExpandable(s *Square, og [][]bool) bool {
	// 右下の1マスが拡張可能かどうかを調べる。
	x := s.Start.X + s.Size
	y := s.Start.Y + s.Size
	ly := len(og)
	lx := len(og[0])
	if x > lx-1 {
		return false
	}
	if y > ly-1 {
		return false
	}
	if og[y][x] {
		return false
	}

	// x方向に拡張可能かどうかを調べる。
	for i := 0; i < s.Size; i++ {
		cy := s.Start.Y + i
		if og[cy][x] {
			return false
		}
	}

	// y方向に拡張可能かどうかを調べる。
	for i := 0; i < s.Size; i++ {
		cx := s.Start.X + i
		if og[y][cx] {
			return false
		}
	}
	return true
}
