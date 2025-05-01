package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// CONFIGURATION
var mode = ANSI

const (
	Gray   = "\033[90m"
	Yellow = "\033[33m"
	Reset  = "\033[0m"
)

type Mode struct {
	pre  string
	post string
}

var (
	ANSI Mode = Mode{Yellow, Reset}
	HTML Mode = Mode{`<span class="poem-diff">`, "</span>"}
)

func readLines(filename string) ([]string, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var lines []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func alignTokens(a, b string) [][2]string {
	aTokens := strings.Fields(a)
	bTokens := strings.Fields(b)

	// Build LCS table
	m, n := len(aTokens), len(bTokens)
	lcs := make([][]int, m+1)
	for i := range lcs {
		lcs[i] = make([]int, n+1)
	}

	for i := m - 1; i >= 0; i-- {
		for j := n - 1; j >= 0; j-- {
			if aTokens[i] == bTokens[j] {
				lcs[i][j] = lcs[i+1][j+1] + 1
			} else {
				if lcs[i+1][j] > lcs[i][j+1] {
					lcs[i][j] = lcs[i+1][j]
				} else {
					lcs[i][j] = lcs[i][j+1]
				}
			}
		}
	}

	// Reconstruct alignment
	var aligned [][2]string
	i, j := 0, 0
	for i < m || j < n {
		if i < m && j < n && aTokens[i] == bTokens[j] {
			aligned = append(aligned, [2]string{aTokens[i], bTokens[j]})
			i++
			j++
		} else if j < n && (i == m || lcs[i][j+1] >= lcs[i+1][j]) {
			aligned = append(aligned, [2]string{"", bTokens[j]})
			j++
		} else if i < m {
			aligned = append(aligned, [2]string{aTokens[i], ""})
			i++
		}
	}

	return aligned
}

func getRow(lineNum int, aLine, bLine string, b Mode) (string, string, string) {
	leftBuilder, rightBuilder := strings.Builder{}, strings.Builder{}

	aligned := alignTokens(aLine, bLine)
	mismatchedLeft, mismatchedRight := strings.Builder{}, strings.Builder{}
	for _, pair := range aligned {
		left, right := pair[0], pair[1]

		if left == right {
			// Flush mismatched tokens before adding matched tokens
			if mismatchedLeft.Len() > 0 {
				leftBuilder.WriteString(b.pre + mismatchedLeft.String() + b.post + " ")
				mismatchedLeft.Reset()
			}
			if mismatchedRight.Len() > 0 {
				rightBuilder.WriteString(b.pre + mismatchedRight.String() + b.post + " ")
				mismatchedRight.Reset()
			}

			leftBuilder.WriteString(left + " ")
			rightBuilder.WriteString(right + " ")
		} else {
			// Accumulate mismatched tokens
			if left != "" {
				if mismatchedLeft.Len() > 0 {
					mismatchedLeft.WriteString(" ")
				}
				mismatchedLeft.WriteString(left)
			}
			if right != "" {
				if mismatchedRight.Len() > 0 {
					mismatchedRight.WriteString(" ")
				}
				mismatchedRight.WriteString(right)
			}
		}
	}

	// Flush any remaining mismatched tokens
	if mismatchedLeft.Len() > 0 {
		leftBuilder.WriteString(b.pre + mismatchedLeft.String() + b.post + " ")
	}
	if mismatchedRight.Len() > 0 {
		rightBuilder.WriteString(b.pre + mismatchedRight.String() + b.post + " ")
	}

	return Gray + strconv.Itoa(lineNum+1) + Reset, leftBuilder.String(), rightBuilder.String()
}

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Usage: go run diff.go file1.txt file2.txt")
		os.Exit(1)
	}

	lines1, err := readLines(os.Args[1])
	if err != nil {
		panic(err)
	}
	lines2, err := readLines(os.Args[2])
	if err != nil {
		panic(err)
	}

	if len(lines1) != len(lines2) {
		panic("unhandled length mismatch")
	}

	numsCol := []string{""}
	leftCol := []string{mode.pre + os.Args[1] + mode.post}
	rightCol := []string{mode.pre + os.Args[2] + mode.post}

	for i := range len(lines1) {
		num, left, right := getRow(i, lines1[i], lines2[i], mode)
		numsCol = append(numsCol, num)
		leftCol = append(leftCol, left)
		rightCol = append(rightCol, right)
	}

	switch mode {
	case ANSI:
		writeAnsi(numsCol, leftCol, rightCol)
		return
	case HTML:
		writeHTML(leftCol, rightCol)
		return
	}
}

func writeAnsi(numsCol, leftCol, rightCol []string) {
	pad(numsCol)
	pad(leftCol)
	pad(rightCol)

	for i := range numsCol {
		fmt.Print(Gray, numsCol[i], Reset, "\t", leftCol[i], "\t", rightCol[i], "\n")
	}
}

func writeHTML(leftCol, rightCol []string) {
	for _, col := range [][]string{leftCol, rightCol} {
		fmt.Print(`<div style="white-space: pre-line">`)
		col = col[1:] // remove header
		for _, s := range col {
			fmt.Println(s)
		}
		fmt.Println(`</div>`)
	}
}

// pad to visual width; see [lipgloss.Width].
func pad(col []string) {
	max := 0
	for _, s := range col {
		if lipgloss.Width(s) > max {
			max = lipgloss.Width(s)
		}
	}

	for i, s := range col {
		padding := max - lipgloss.Width(s)
		if padding > 0 {
			col[i] += strings.Repeat(" ", padding)
		}
	}
}
