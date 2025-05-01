package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	Gray   = "\033[90m"
	Yellow = "\033[33m"
	Reset  = "\033[0m"
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

func getRow(lineNum int, aLine, bLine string) (string, string, string) {
	leftBuilder, rightBuilder := strings.Builder{}, strings.Builder{}

	aligned := alignTokens(aLine, bLine)
	for _, pair := range aligned {
		left, right := pair[0], pair[1]

		if left == right {
			leftBuilder.WriteString(left + " ")
			rightBuilder.WriteString(right + " ")
		} else {
			if left != "" {
				leftBuilder.WriteString(Yellow + left + Reset + " ")
			}
			if right != "" {
				rightBuilder.WriteString(Yellow + right + Reset + " ")
			}
		}
	}

	return Gray + strconv.Itoa(lineNum+1) + Reset, leftBuilder.String(), rightBuilder.String()
}

func main() {
	if len(os.Args) != 3 {
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
	leftCol := []string{Gray + os.Args[1] + Reset}
	rightCol := []string{Gray + os.Args[2] + Reset}

	for i := range len(lines1) {
		num, left, right := getRow(i, lines1[i], lines2[i])
		numsCol = append(numsCol, num)
		leftCol = append(leftCol, left)
		rightCol = append(rightCol, right)
	}

	pad(numsCol)
	pad(leftCol)
	pad(rightCol)

	for i := range numsCol {
		fmt.Print(Gray, numsCol[i], Reset, "\t", leftCol[i], "\t", rightCol[i], "\n")
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
