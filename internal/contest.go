package internal

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"github.com/hunderaweke/codative-codeforces/session"
)

type Contest struct {
	Problems    []problem
	Title       string
	ContestType string
	ContestID   string
}

var (
	wg sync.WaitGroup
	mu sync.Mutex
)

func (c *Contest) Create(directoryName string, template Template) error {
	err := os.Chdir(c.ContestType)
	if err != nil {
		os.Mkdir(c.ContestType, 0777)
		os.Chdir(c.ContestType)
	}
	os.Mkdir(directoryName, 0777)
	if err = os.Chdir(directoryName); err != nil {
		return (err)
	}
	data, err := template.Load()
	ext := FileExtensions[template.Lang]
	if err != nil {
		return err
	}
	for _, prob := range c.Problems {
		if err := prob.create(data, ext); err != nil {
			return err
		}
	}
	wg.Wait()
	return nil
}

func Parse(contestID, contestType string) Contest {
	c := Contest{ContestID: contestID, ContestType: contestType, Problems: []problem{}}
	c.Title = getContestTitle(contestID, contestType)
	probs := findProblems(contestID, contestType)
	for id := range probs {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			p := parseProblem(contestID, contestType, id)
			p.title = id + ". " + probs[id]
			mu.Lock()
			c.Problems = append(c.Problems, p)
			mu.Unlock()
		}(id)
	}
	wg.Wait()
	return c
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
