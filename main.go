package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func main() {
	// 引数チェック（ファイルパスが貼り付けられてくる想定）
	if len(os.Args) < 2 {
		fmt.Println("使用法: go run main.go <ファイルパス>")
		return
	}
	filePath := os.Args[1]

	file, err := os.Open(filePath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "エラー: %v\n", err)
		os.Exit(1)
	}
	defer file.Close()

	process(file, os.Stdout)
}

func process(r io.Reader, w io.Writer) {
	reader := bufio.NewReader(r)
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	inString := false
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}

		if r == '"' {
			inString = !inString
			writer.WriteRune(r)
			continue
		}

		if !inString && r == '/' {
			next, _, _ := reader.ReadRune()
			if next == '/' { // 1行コメント
				for {
					n, _, _ := reader.ReadRune()
					if n == '\n' || n == 0 {
						writer.WriteRune('\n')
						break
					}
				}
				continue
			} else if next == '*' { // 複数行コメント
				for {
					n, _, _ := reader.ReadRune()
					if n == '*' {
						nn, _, _ := reader.ReadRune()
						if nn == '/' {
							break
						}
					}
					if n == 0 { break }
				}
				continue
			} else {
				writer.WriteRune(r)
				_ = reader.UnreadRune()
				continue
			}
		}
		writer.WriteRune(r)
	}
}
