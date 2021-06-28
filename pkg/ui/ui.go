package ui

import (
	"embed"
	"io/fs"
)

//go:embed html
var html embed.FS

var htmlFolder = "html"

type HtmlPage struct {
	htmlFS embed.FS
}

func NewHtmlPage() *HtmlPage {
	return &HtmlPage{
		htmlFS: html,
	}
}

func (w *HtmlPage) Get() embed.FS {
	return w.htmlFS
}

func (w *HtmlPage) GetFS() fs.FS {
	files, _ := fs.Sub(w.htmlFS, htmlFolder)
	return files
}

func (w *HtmlPage) GetTemplate(filename string) (string, error) {
	var content []byte
	err := fs.WalkDir(w.GetFS(), ".", func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if !d.IsDir() {
			if s == filename {
				b, err := w.htmlFS.ReadFile(htmlFolder + "/" + s)
				if err != nil {
					return err
				}
				content = b
			}
		}
		return nil
	})
	return string(content), err
}
