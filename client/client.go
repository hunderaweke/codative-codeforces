package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"os"

	"github.com/fatih/color"
	"github.com/juju/persistent-cookiejar"
)

type Client struct {
	Jar           *cookiejar.Jar `json:"cookies,omitempty"`
	Handle        string         `json:"handle,omitempty"`
	HandleOrEmail string         `json:"handle_or_email,omitempty"`
	Password      string         `json:"password,omitempty"`
	Ftaa          string         `json:"ftaa,omitempty"`
	Bfaa          string         `json:"bfaa,omitempty"`
	host          string
	path          string
	client        *http.Client
}

var Clnt *Client

func Create(host, path string) {
	jar, _ := cookiejar.New(nil)
	c := &Client{Jar: jar, host: host, path: path, client: nil}
	if err := c.load(); err != nil {
		color.Red("%s", "Session file not found")
		color.Blue("%s", "Creating new configuration file")
	}
	c.client = &http.Client{Jar: c.Jar}
	if err := c.save(); err != nil {
		fmt.Println(err)
		color.Red("%s", "Cannot Save the configuration")
	}
	Clnt = c
}

func (c *Client) load() error {
	os.Chdir(c.path)
	file, err := os.Open(".codative.session")
	if err != nil {
		return err
	}

	data, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	defer file.Close()
	json.Unmarshal(data, c)
	return nil
}

func (c *Client) save() error {

	os.Chdir(c.path)
	file, err := os.Create(".codative.session")
	if err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", " ")
	if err != nil {
		return err
	}

	defer file.Close()
	_, err = file.Write(data)
	if err != nil {
		return err
	}
	return nil
}
