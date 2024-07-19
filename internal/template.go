package internal

import (
	"os"

	"github.com/fatih/color"
)

type Template struct {
	Lang  string `json:"lang,omitempty"`
	Path  string `json:"path,omitempty"`
	Alias string `json:"alias"`
}

func (t *Template) Load() ([]byte, error) {
	bytes, err := os.ReadFile(t.Path)
	if err != nil {
		color.Red("Error Reading template file %v\n", err)
		return nil, err
	}
	return bytes, nil
}

func (t *Template) createFile(fileName string) error {
	fileExtension := FileExtensions[t.Lang]
	fileName += fileExtension
	os.Create(fileName)
	templateContent, err := t.Load()
	if err != nil {
		return err
	}
	os.WriteFile(fileName, templateContent, 0644)
	return nil
}

func NewTemplate(filePath, lang, alias string) Template {
	t := Template{Path: filePath, Lang: lang, Alias: alias}
	return t
}
