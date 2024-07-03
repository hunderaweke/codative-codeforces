package utils

import (
	"fmt"
	"net/http"
	"os"

	"github.com/PuerkitoBio/goquery"
)

func FetchTestCases(contestID int, index string, probName string) {
	url := fmt.Sprintf("https://codeforces.com/contest/%d/problem/%s", contestID, index)
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching problem page: %v\n", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Printf("Failed to fetch page with status code: %d\n", resp.StatusCode)
		return
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		fmt.Printf("Error parsing HTML: %v\n", err)
		return
	}

	doc.Find(".input pre").Each(func(i int, s *goquery.Selection) {
		s.Find("br").Each(func(_ int, br *goquery.Selection) {
			br.ReplaceWithHtml("\n")
		})
		fileName := fmt.Sprintf("input%d.in", i+1)
		os.Create(fileName)
		os.WriteFile(fileName, []byte(s.Text()), 0755)
	})

	doc.Find(".output pre").Each(func(i int, s *goquery.Selection) {
		s.Find("br").Each(func(_ int, br *goquery.Selection) {
			br.ReplaceWithHtml("\n")
		})
		fileName := fmt.Sprintf("output%d.out", i+1)
		os.Create(fileName)
		os.WriteFile(fileName, []byte(s.Text()), 0755)
	})
}
