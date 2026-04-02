package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("使い方: cs-strip.exe <ファイルパス>")
		return
	}
	inputPath := os.Args[1]
	// パスの前後にある引用符（"）を削除（コピペ対策）
	inputPath = strings.Trim(inputPath, "\"")

	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("ファイルが開けません: %v\n", err)
		return
	}
	defer file.Close()

	outputPath := inputPath + "_cleaned.cs"
	outFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("出力ファイルが作成できません: %v\n", err)
		return
	}
	defer outFile.Close()

	process(file, outFile)
	fmt.Printf("完了！保存先: %s\n", outputPath)
}

func process(r io.Reader, w io.Writer) {
	reader := bufio.NewReader(r)
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	inString := false
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF { break }

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
						if nn == '/' { break }
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
