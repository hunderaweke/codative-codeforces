package config

import (
	"encoding/json"
	"os"

	"github.com/fatih/color"
	"github.com/hunderaweke/codative-codeforces/internal"
)

type LocalConfig struct {
	Template  internal.Template `json:"template"`
	ContestID string            `json:"contest_id"`
}

func (c *LocalConfig) Save(path string) error {
	err := os.Chdir(path)
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "    ")
	if err != nil {
		return err
	}
	file, err := os.Create(".codative.json")
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}

func (c *LocalConfig) Load(path string) error {
	err := os.Chdir(path)
	if err != nil {
		return err
	}
	file, err := os.ReadFile(".codative.json")
	if err != nil {
		return err
	}
	json.Unmarshal(file, c)
	return nil
}

func LocalConfigPrompt() {
	var globalConfig GlobalConfig
	globalConfig.ConfigPath = "/home/hundera/.codative/"
	if err := globalConfig.Load(); err != nil {
		color.Red("Failed to load the global configutation file %v", err.Error())
	}

}
