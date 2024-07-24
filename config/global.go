package config

import (
	"encoding/json"
	"os"

	"github.com/hunderaweke/codative-codeforces/internal"
)

type GlobalConfig struct {
	BaseDir    string              `json:"base_dir"`
	Handle     string              `json:"handle,omitempty"`
	Host       string              `json:"host,omitempty"`
	Templates  []internal.Template `json:"templates,omitempty"`
	configPath string
}

func (c *GlobalConfig) Load() error {
	err := os.Chdir(c.configPath)
	if err != nil {
		return err
	}
	file, err := os.ReadFile("config.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, c)
	if err != nil {
		return err
	}
	return nil
}

func (c *GlobalConfig) Save() error {
	err := os.Chdir(c.configPath)
	if err != nil {
		return err
	}
	file, err := os.Create("config.json")
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}
	file.Write(data)
	return nil
}

func (c *GlobalConfig) ConfigPrompt() {
}
