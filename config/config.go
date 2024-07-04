package config

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/fatih/color"
)

type CodeTemplate struct {
	Alias        string `json:"alias,omitempty"`
	Lang         string `json:"lang,omitempty"`
	Path         string `json:"path,omitempty"`
	Suffix       string `json:"suffix,omitempty"`
	BeforeScript string `json:"before_script,omitempty"`
	Script       string `json:"script,omitempty"`
	AfterScript  string `json:"after_script,omitempty"`
}
type Config struct {
	Templates       []CodeTemplate    `json:"templates,omitempty"`
	DefaultTemplate int               `json:"default_template,omitempty"`
	GenAfterParse   bool              `json:"gen_after_parse,omitempty"`
	Host            string            `json:"host,omitempty"`
	Proxy           string            `json:"proxy,omitempty"`
	FolderName      map[string]string `json:"folder_name,omitempty"`
	path            string
}

var Instance *Config

func Init(path string) {
	c := &Config{path: path, Host: "https://codeforces.com", Proxy: ""}
	if err := c.load(); err != nil {
		color.Green("Creating a new configuration in %v \n", path)
	}
	if c.DefaultTemplate < 0 || c.DefaultTemplate >= len(c.Templates) {
		c.DefaultTemplate = 0
	}
	if c.FolderName == nil {
		c.FolderName = map[string]string{}
	}
	if _, ok := c.FolderName["root"]; !ok {
		c.FolderName["root"] = ".codative_config"
	}
	c.save()
	Instance = c

}

func (c *Config) load() (err error) {
	file, err := os.Open(c.path)
	if err != nil {
		return
	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return
	}
	return json.Unmarshal(bytes, c)
}

func (c *Config) save() (err error) {
	var data bytes.Buffer
	encoder := json.NewEncoder(&data)
	encoder.SetIndent("", "  ")
	encoder.SetEscapeHTML(false)
	err = encoder.Encode(c)
	if err == nil {
		os.MkdirAll(filepath.Dir(c.path), os.ModePerm)
		err = ioutil.WriteFile(c.path, data.Bytes(), 0644)
	}
	if err != nil {
		color.Red("Cannot Save file in the directory %v\n%v", c.path, err.Error())
	}
	return
}
