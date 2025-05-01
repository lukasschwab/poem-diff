package main

import (
	_ "embed"
	"os"
	"strings"

	"fmt"
	"html/template"
	"io"
)

func writeHTML(leftCol, rightCol []string) {
	entries := make([]htmlEntry, 2)
	entries[0].Title = leftCol[0]
	entries[1].Title = rightCol[0]
	for i, col := range [][]string{leftCol, rightCol} {
		var b strings.Builder
		b.WriteString(`<div style="white-space: pre-line">`)

		col = col[1:] // remove header
		for _, s := range col {
			fmt.Fprintln(&b, s)
		}
		b.WriteString(`</div>`)

		entries[i].Text = template.HTML(b.String())
	}
	renderHTML(entries, os.Stdout)
}

//go:embed templates/html.tmpl
var htmlTemplate string

type htmlEntry struct {
	Title string
	Text  template.HTML
}

func renderHTML(entries []htmlEntry, to io.Writer) error {
	tmpl, err := template.New("html").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("error parsing HTML template: %w", err)
	}

	return tmpl.Execute(to, entries)
}
