package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	// 最後に必ず入力を待つ
	defer func() {
		fmt.Println("\n-----------------------------------------")
		fmt.Println("処理が終了しました。Enterキーを押すと閉じます...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}()

	// 引数チェック
	if len(os.Args) < 2 {
		fmt.Println("【使い方】")
		fmt.Println("このexeファイルにC#ファイルをドラッグ＆ドロップするか、")
		fmt.Println("パスを貼り付けてEnterを押してください。")
		return
	}

	// 渡されたパスを整理
	inputPath := strings.Trim(os.Args[1], "\"")
	fmt.Printf("読み込み中: %s\n", inputPath)

	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("エラー: ファイルが開けませんでした: %v\n", err)
		return
	}
	defer file.Close()

	// 出力ファイル名を作成
	outputPath := inputPath + "_cleaned.cs"
	outFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("エラー: 出力ファイルが作成できませんでした: %v\n", err)
		return
	}
	defer outFile.Close()

	stripComments(file, outFile)

	fmt.Printf("成功！保存先:\n%s\n", outputPath)
}

func stripComments(r io.Reader, w io.Writer) {
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
			next, _, err := reader.ReadRune()
			if err != nil {
				writer.WriteRune(r)
				break
			}

			if next == '/' { // 1行コメント
				for {
					n, _, err := reader.ReadRune()
					if err == io.EOF || n == '\n' {
						writer.WriteRune('\n')
						break
					}
				}
				continue
			} else if next == '*' { // 複数行コメント
				for {
					n, _, err := reader.ReadRune()
					if err == io.EOF {
						break
					}
					if n == '*' {
						nn, _, _ := reader.ReadRune()
						if nn == '/' {
							break
						}
						_ = reader.UnreadRune()
					}
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
