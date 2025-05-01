package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
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

// TODO: slop?
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

func printAlignedTokens(lineNum int, aLine, bLine string) {
	fmt.Printf("%d:\t", lineNum+1)

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

	fmt.Print(leftBuilder.String(), "\t", rightBuilder.String(), "\n")
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

	longestLeft := 0
	for _, line := range lines1 {
		if longestLeft < len(line) {
			longestLeft = len(line)
		}
	}

	maxLines := len(lines1)
	if len(lines2) > maxLines {
		maxLines = len(lines2)
	}

	// TODO: figure out padding
	fmt.Print("\t", os.Args[1], "\t", os.Args[2], "\n")

	for i := 0; i < maxLines; i++ {
		var l1, l2 string
		if i < len(lines1) {
			l1 = lines1[i]
		}
		if i < len(lines2) {
			l2 = lines2[i]
		}

		if l1 == l2 {
			continue
		}

		printAlignedTokens(i, l1, l2)
	}
}
