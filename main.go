package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
)

func main() {
	myApp := app.New()
	myWindow := myApp.NewWindow("C# Comment Stripper")
	myWindow.Resize(fyne.NewSize(400, 200))

	label := widget.NewLabel("C#ファイルをここにドロップするか\nボタンで選択してください")
	label.Alignment = fyne.TextAlignCenter

	// 処理関数
	runStrip := func(path string) {
		err := processFile(path)
		if err != nil {
			dialog.ShowError(err, myWindow)
		} else {
			dialog.ShowInformation("完了", "コメントを削除したファイルを保存しました:\n"+path+"_cleaned.cs", myWindow)
		}
	}

	// ボタン
	btn := widget.NewButton("ファイルを選択", func() {
		dialog.ShowFileOpen(func(reader fyne.URIReadCloser, err error) {
			if err == nil && reader != nil {
				runStrip(reader.URI().Path())
			}
		}, myWindow)
	})

	// ドラッグ＆ドロップ対応
	myWindow.SetOnDropped(func(pos fyne.Position, uris []fyne.URI) {
		for _, uri := range uris {
			runStrip(uri.Path())
		}
	})

	myWindow.SetContent(container.NewVBox(
		label,
		btn,
	))

	myWindow.ShowAndRun()
}

func processFile(inputPath string) error {
	inputPath = strings.Trim(inputPath, "\"")
	file, err := os.Open(inputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	outputPath := inputPath + "_cleaned.cs"
	outFile, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	stripComments(file, outFile)
	return nil
}

func stripComments(r io.Reader, w io.Writer) {
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
			next, _, err := reader.ReadRune()
			if err != nil {
				writer.WriteRune(r)
				break
			}
			if next == '/' {
				for {
					n, _, err := reader.ReadRune()
					if err == io.EOF || n == '\n' {
						writer.WriteRune('\n')
						break
					}
				}
				continue
			} else if next == '*' {
				for {
					n, _, err := reader.ReadRune()
					if err == io.EOF { break }
					if n == '*' {
						nn, _, _ := reader.ReadRune()
						if nn == '/' { break }
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
