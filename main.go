package main

import (
	"fmt"
	"os"
)

func main() {
	// プログラム引数で複数のマップファイルが指定されたときの処理を行う。
	args := os.Args
	files := args[1:]
	for _, file := range files {
		func() {
			f, err := os.Open(file)
			if err != nil {
				fmt.Fprintln(os.Stderr, err)
				return
			}
			defer f.Close()
			Bsq(f, os.Stdout, os.Stderr)
		}()
	}

	// プログラム引数がない場合は、標準入力からマップを読み込んで処理を行う。
	l := len(files)
	if l == 0 {
		Bsq(os.Stdin, os.Stdout, os.Stderr)
	}
}
