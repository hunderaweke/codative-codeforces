package utils

import (
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"unicode"
)

func PostBody(client *http.Client, u *url.URL, data url.Values) ([]byte, error) {
	resp, err := client.PostForm(u.String(), data)
	if err != nil {
		return nil, err
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func GetBody(client *http.Client, u *url.URL) ([]byte, error) {
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// CHA map
const CHA = "abcdefghijklmnopqrstuvwxyz0123456789"

// RandString n is the length. a-z 0-9
func RandString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = CHA[rand.Intn(len(CHA))]
	}
	return string(b)
}

func ReformString(name string) string {
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
