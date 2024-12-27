package helpers

import (
	"bufio"
	"io"
)

func ReadLine(r io.Reader, lineNum int) (string, error) {
	lastLine := 0
	sc := bufio.NewScanner(r)

	for sc.Scan() {
		lastLine++

		if lastLine == lineNum {
			return sc.Text(), sc.Err()
		}
	}

	return "", io.EOF
}
