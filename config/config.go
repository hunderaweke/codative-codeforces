package config

import (
	"encoding/json"
	"os"

	"github.com/mitchellh/go-homedir"
)

type Config interface {
	save() error
}
type LocalConfig struct {
	Lang       string `json:"lang,omitempty"`
	TemplateId int    `json:"template_id,omitempty"`
	ContestId  string `json:"contest_id,omitempty"`
}
type GlobalConfig struct {
	Handle    string     `json:"handle,omitempty"`
	Host      string     `json:"host,omitempty"`
	Templates []Template `json:"templates,omitempty"`
}

func load() (Config, error) {
	//TODO: Loading both Global and Local Config

}

func (c *LocalConfig) save() error {
	//TODO: Saving the local config to the Contest path
	return nil
}
func (c *GlobalConfig) save() error {
	homeDir, err := homedir.Dir()
	if err != nil {
		return err
	}
	_, err = os.ReadFile(homeDir + ".codative_config")
	if err == nil {
		os.Remove(homeDir + ".codative_config")
	}
	file, err := os.Create(homeDir + ".codative_config")
	bytes, err := json.Marshal(c)
	file.Write(bytes)
	return nil
}
