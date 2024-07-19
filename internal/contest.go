package internal

import (
	"fmt"
	"os"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/hunderaweke/codative-codeforces/config"
	"github.com/hunderaweke/codative-codeforces/session"
	"github.com/hunderaweke/codative-codeforces/utils"
)

type Contest struct {
	Problems    []Problem
	Title       string
	ContestType string
}

func (c *Contest) Create() error {
	// TODO: Implement creation file for specific contest
	err := os.Chdir(strings.Join([]string{config.C.BaseDir, c.ContestType}, "/"))
	if err != nil {
		return err
	}
	contestDir := utils.ReformString(c.Title)
	os.Mkdir(contestDir, 0644)
	return nil
}

func Parse(contestID, contestType string) error {
	// TODO: Implement the parsing of Contest and creation of files for the contest
}

func getContestTitle(contestID, contestType string) string {
	r, _ := session.S.Client.Get(session.S.Host + "/" + contestType + "/" + contestID)
	doc, _ := goquery.NewDocumentFromReader(r.Body)
	s := doc.Find(".rtable")
	ret := ""
	s.Find("th").Each(func(i int, s *goquery.Selection) {
		ret = fmt.Sprintf("%v", s.Has("a").Text())
	})
	return ret
}

func findProblems(contestID, contestType string) map[string]string {
	r, _ := session.S.Client.Get(session.S.Host + contestType + "/" + contestID)
	doc, _ := goquery.NewDocumentFromReader(r.Body)
	defer r.Body.Close()
	s := doc.Find(".problems").Find("tr")
	problems := make(map[string]string)
	key := ""
	s.Find("tr > td").Each(func(i int, s *goquery.Selection) {
		index := strings.TrimSpace(s.Find(".id > a").Text())
		if len(index) > 0 {
			key = index
		}
		name := s.Find("div>a").Text()
		if len(name) > 0 {
			problems[key] = name
			key = ""
		}
	})
	return problems
}
