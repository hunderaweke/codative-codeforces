package session

import (
	"encoding/json"
	"net/http"
	"net/url"
	"os"

	persistent "github.com/juju/persistent-cookiejar"
	"golang.org/x/net/publicsuffix"
)

type Session struct {
	Cookies       []*http.Cookie `json:"cookies"`
	Ftaa          string         `json:"ftaa,omitempty"`
	Bfaa          string         `json:"bfaa,omitempty"`
	Host          string         `json:"host"`
	HandleOrEmail string         `json:"handle_or_email"`
	Handle        string         `json:"handle"`
	Password      string         `json:"password"`
	path          string         `json:"-"`
	Client        *http.Client   `json:"-"`
}

var S *Session

func init() {
	s := &Session{Cookies: []*http.Cookie{}, Host: "https://codeforces.com/", Client: &http.Client{}}
	home, _ := os.UserHomeDir()
	s.path = home + "/.codative"
	S = s
}

func (s *Session) Load() error {
	err := os.Chdir(s.path)
	if err != nil {
		return err
	}
	file, err := os.ReadFile("session.json")
	if err != nil {
		return err
	}
	err = json.Unmarshal(file, s)
	if err != nil {
		return err
	}

	jar, _ := persistent.New(&persistent.Options{PublicSuffixList: publicsuffix.List})
	u, _ := url.Parse(S.Host)
	jar.SetCookies(u, S.Cookies)
	S.Client = &http.Client{Jar: jar}
	return nil
}

func (s *Session) Save() error {
	err := os.Chdir(s.path)
	if err != nil {
		return err
	}
	file, err := os.Create("session.json")
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := json.MarshalIndent(s, "", " ")
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	return nil
}
