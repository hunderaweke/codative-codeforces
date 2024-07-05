package config

import (
	"os"

	"github.com/fatih/color"
	"github.com/hunderaweke/codative-codeforces/client"
)

type Template struct {
	Lang string `json:"lang,omitempty"`
	Path string `json:"path,omitempty"`
	Name string `json:"name,omitempty"`
}

func (t *Template) load() ([]byte, error) {

	bytes, err := os.ReadFile(t.Path)
	if err != nil {
		color.Red("Error Reading template file %v\n", err)
		return nil, err
	}
	return bytes, nil
}

func (t *Template) createFile(fileName string) error {
	fileExtension := client.FILE_EXTENSIONS[client.LANGS[t.Lang]]
	fileName += fileExtension
	os.Create(fileName)
	templateContent, err := t.load()
	if err != nil {
		return err
	}
	os.WriteFile(fileName, templateContent, 0644)
	return nil
}
