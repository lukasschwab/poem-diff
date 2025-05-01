package main

import (
	"bufio"
	"flag"
	"os"
	"strconv"
	"strings"
)

const (
	ANSIGray   = "\033[90m"
	ANSIYellow = "\033[33m"
	ANSIReset  = "\033[0m"
)

type Mode struct {
	pre  string
	post string
}

var (
	ANSI Mode = Mode{ANSIYellow, ANSIReset}
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

// WARNING: this is a LLM-generated implementation of the LCS algorithm at the
// token level. It might suck, but I'm only dealing with ~hundreds of tokens.
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

// WARNING: this is LLM-generated code.
//
// getRow returns the line number and formats the left/right texts, inserting
// markers (specified by mode) before and after serieses of mismatched tokens.
//
// It may be unnecessarily complex, since those spans cover different sets of
// tokens on the left and right sides. The space-appending is weird to avoid
// including spaces in those mismatch spans.
func getRow(lineNum int, aLine, bLine string, mode Mode) (string, string, string) {
	leftBuilder, rightBuilder := strings.Builder{}, strings.Builder{}

	aligned := alignTokens(aLine, bLine)
	mismatchedLeft, mismatchedRight := strings.Builder{}, strings.Builder{}
	for _, pair := range aligned {
		left, right := pair[0], pair[1]

		if left == right {
			// Flush mismatched tokens before adding matched tokens
			if mismatchedLeft.Len() > 0 {
				leftBuilder.WriteString(mode.pre + mismatchedLeft.String() + mode.post + " ")
				mismatchedLeft.Reset()
			}
			if mismatchedRight.Len() > 0 {
				rightBuilder.WriteString(mode.pre + mismatchedRight.String() + mode.post + " ")
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
		leftBuilder.WriteString(mode.pre + mismatchedLeft.String() + mode.post + " ")
	}
	if mismatchedRight.Len() > 0 {
		rightBuilder.WriteString(mode.pre + mismatchedRight.String() + mode.post + " ")
	}

	return ANSIGray + strconv.Itoa(lineNum+1) + ANSIReset, leftBuilder.String(), rightBuilder.String()
}

func main() {
	left := flag.String("left", "", "Path to the left file")
	right := flag.String("right", "", "Path to the right file")
	format := flag.String("format", "ansi", "Output format: 'ansi' or 'html'")
	flag.Parse()

	if left == nil || right == nil || *left == "" || *right == "" {
		panic("specify files --left and --right")
	}

	lines1, err := readLines(*left)
	if err != nil {
		panic(err)
	}
	lines2, err := readLines(*right)
	if err != nil {
		panic(err)
	}

	if len(lines1) != len(lines2) {
		panic("unhandled length mismatch")
	}

	mode := ANSI
	if format != nil {
		switch *format {
		case "ansi":
			mode = ANSI
		case "html":
			mode = HTML
		default:
			panic("Invalid format. Use 'ansi' or 'html'")
		}
	}

	numsCol := []string{""}
	leftCol := []string{*left}
	rightCol := []string{*right}
	for i := range lines1 {
		num, left, right := getRow(i, lines1[i], lines2[i], mode)
		numsCol = append(numsCol, num)
		leftCol = append(leftCol, left)
		rightCol = append(rightCol, right)
	}

	switch mode {
	case ANSI:
		writeAnsi(numsCol, leftCol, rightCol)
	case HTML:
		writeHTML(leftCol, rightCol)
	}
}
