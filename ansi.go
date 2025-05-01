package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func writeAnsi(numsCol, leftCol, rightCol []string) {
	pad(numsCol)
	pad(leftCol)
	pad(rightCol)

	fmt.Println("")
	for i := range numsCol {
		fmt.Print(ANSIGray, numsCol[i], ANSIReset, "\t", leftCol[i], "\t", rightCol[i], "\n")
		if i == 0 {
			fmt.Println("")
		}
	}
	fmt.Println("")
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
