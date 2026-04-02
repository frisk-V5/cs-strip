package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func main() {
	// 最後に必ず入力を待つように設定
	defer func() {
		fmt.Println("\n-----------------------------------------")
		fmt.Println("処理が終了しました。Enterキーを押すと閉じます...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}()

	// 引数チェック（ファイルパスが渡されているか）
	if len(os.Args) < 2 {
		fmt.Println("【使い方】")
		fmt.Println("このexeファイルにC#ファイルをドラッグ＆ドロップするか、")
		fmt.Println("コマンドプロンプトで 'cs-strip.exe ファイルパス' と入力してください。")
		return
	}

	// 渡されたパスを整理（前後にある引用符 \" を削除）
	inputPath := strings.Trim(os.Args[1], "\"")
	fmt.Printf("読み込み中: %s\n", inputPath)

	// 入力ファイルを開く
	file, err := os.Open(inputPath)
	if err != nil {
		fmt.Printf("エラー: ファイルが開けませんでした: %v\n", err)
		return
	}
	defer file.Close()

	// 出力ファイル名を作成 (例: Program.cs -> Program_cleaned.cs)
	outputPath := inputPath + "_cleaned.cs"
	outFile, err := os.Create(outputPath)
	if err != nil {
		fmt.Printf("エラー: 出力ファイルが作成できませんでした: %v\n", err)
		return
	}
	defer outFile.Close()

	// コメント削除処理の実行
	stripComments(file, outFile)

	fmt.Printf("成功！コメントを削除したファイルを保存しました:\n%s\n", outputPath)
}

// コメント削除のメインロジック
func stripComments(r io.Reader, w io.Writer) {
	reader := bufio.NewReader(r)
	writer := bufio.NewWriter(w)
	defer writer.Flush()

	inString := false // 文字列リテラル中かどうか
	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			break
		}

		// 文字列リテラル（"..."）の開始・終了判定
		if r == '"' {
			inString = !inString
			writer.WriteRune(r)
			continue
		}

		// 文字列外で '/' が出たらコメントの可能性あり
		if !inString && r == '/' {
			next, _, _ := reader.ReadRune()
			if next == '/' { // 1行コメント // ...
				for {
					n, _, _ := reader.ReadRune()
					if n == '\n' || n == 0 {
						writer.WriteRune('\n')
						break
					}
				}
				continue
			} else if next == '*' { // 複数行コメント /* ... */
				for {
					n, _, _ := reader.NewReader(reader).ReadRune() // 実際には直読み
					n, _, _ = reader.ReadRune()
					if n == '*' {
						nn, _, _ := reader.ReadRune()
						if nn == '/' {
							break
						}
					}
					if n == 0 {
						break
					}
				}
				continue
			} else {
				// コメントではなかったので '/' と次の文字をそのまま出力
				writer.WriteRune(r)
				_ = reader.UnreadRune()
				continue
			}
		}

		// それ以外の文字はそのまま出力
		writer.WriteRune(r)
	}
}
