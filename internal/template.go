package internal

import (
	"os"
	"strings"

	"github.com/AlecAivazis/survey/v2"
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
func TemplatePrompt() (Template, error) {
	var t Template
	var options []string
	for id, lang := range Langs {
		options = append(options, id+" "+lang)
	}
	q := []*survey.Question{
		{
			Name: "path",
			Prompt: &survey.Input{
				Message: "Enter the template path (absolute path):",
			},
			Validate: survey.Required,
		},
		{
			Name: "lang",
			Prompt: &survey.Select{
				Message: "Choose Language",
				Options: options,
				Help:    "This is the language that will be used while submitting",
			},
			Validate: survey.Required,
		},
		{
			Name: "alias",
			Prompt: &survey.Input{
				Message: "Insert an alias for the template",
				Help:    "This is the alias that will be used for loading and searching the alias",
			},
			Validate: survey.Required,
		},
	}
	survey.Ask(q, &t)
	t.Lang = strings.Split(t.Lang, " ")[0]
	return t, nil
}
