package main

import (
	_ "embed"

	"fmt"
	"html/template"
	"io"
)

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
