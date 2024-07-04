package client

import "net/http/cookiejar"

type Client struct {
	Jar            *cookiejar.Jar `json:"jar,omitempty"`
	Handle         string         `json:"handle,omitempty"`
	HandleOrEmail  string         `json:"handle_or_email,omitempty"`
	Password       string         `json:"password,omitempty"`
	Ftaa           string         `json:"ftaa,omitempty"`
	Bfaa           string         `json:"bfaa,omitempty"`
	LastSubmission string         `json:"last_submission,omitempty"`
	host           string
	proxy          string
}
