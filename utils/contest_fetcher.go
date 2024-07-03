package utils

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/hunderaweke/codative-codeforces/types"
)

func Fetch() (*types.Response, error) {
	resp, err := http.Get("https://codeforces.com/api/contest.standings?contestId=1986")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var r types.Response
	json.Unmarshal(body, &r)
	return &r, nil
}
