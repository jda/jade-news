package main

import (
	"html/template"
	"io"
)

// TemplateData is container for data we pass to remplate
type TemplateData struct {
	Links          []Link
	LatestNewsTime string
}

func renderLinks(w io.Writer, links []Link, tmplFileName string) error {
	t, err := template.ParseFiles(tmplFileName)
	if err != nil {
		return err
	}

	latestNews := links[0].Published.Format("on 02 January 2006 at 15:04 MST")

	td := TemplateData{
		Links:          links,
		LatestNewsTime: latestNews,
	}

	if err := t.ExecuteTemplate(w, tmplFileName, td); err != nil {
		return err
	}

	return nil
}
