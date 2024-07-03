package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"unicode"
)

type Response struct {
	Status string `json:"status,omitempty"`
	Result Result `json:"result,omitempty"`
}
type Result struct {
	Contest  Contest   `json:"contest,omitempty"`
	Problems []Problem `json:"problems,omitempty"`
}

type Contest struct {
	Id     int    `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Phase  string `json:"phase,omitempty"`
	Frozen bool   `json:"frozen,omitempty"`
}
type Problem struct {
	Index string `json:"index,omitempty"`
	Name  string `json:"name,omitempty"`
}

func reformString(name string) string {
	var reformedString []rune
	for _, ch := range name {
		if unicode.IsSymbol(ch) || unicode.IsPunct(ch) {
			continue
		}
		if unicode.IsSpace(ch) {
			reformedString = append(reformedString, '_')
			continue
		}
		reformedString = append(reformedString, ch)
	}
	return string(reformedString)
}

func main() {
	resp, err := http.Get("https://codeforces.com/api/contest.standings?contestId=566")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	var decoded Response
	json.Unmarshal(body, &decoded)
	directoryName := reformString(decoded.Result.Contest.Name)

	os.Mkdir(string(directoryName), 0755)
	os.Chdir(string(directoryName))
	for _, prob := range decoded.Result.Problems {
		probName := reformString(prob.Index + " " + prob.Name)
		os.Create(probName + ".py")
	}
}
