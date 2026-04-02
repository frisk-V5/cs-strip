package main

import (
	"bufio"
	"io"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	writer := bufio.NewWriter(os.Stdout)
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
			if next == '/' { // 1行コメント //
				for {
					n, _, _ := reader.ReadRune()
					if n == '\n' || n == 0 {
						writer.WriteRune('\n')
						break
					}
				}
				continue
			} else if next == '*' { // 複数行コメント /* */
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
				reader.UnreadRune()
				continue
			}
		}
		writer.WriteRune(r)
	}
}
