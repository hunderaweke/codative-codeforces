package session

import (
	"fmt"
	"regexp"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
)

type Problem struct {
	statement string
	input     []string
	output    []string
}

type Contest struct {
	Problems    []Problem
	Title       string
	ContestType string
}

func GetContestTitle(contestID, contestType string) string {
	r, _ := S.Client.Get(S.Host + "/" + contestType + "/" + contestID)
	doc, _ := goquery.NewDocumentFromReader(r.Body)
	s := doc.Find(".rtable")
	ret := ""
	s.Find("th").Each(func(i int, s *goquery.Selection) {
		ret = fmt.Sprintf("%v", s.Has("a").Text())
	})
	return ret
}

func FindProblems(contestID, contestType string) map[string]string {
	r, _ := S.Client.Get(S.Host + contestType + "/" + contestID)
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

func ParseProblem(contestID, contestType, problemID string) (p Problem) {
	u := strings.Join([]string{S.Host, contestType, contestID, "problem", problemID}, "/")
	r, _ := S.Client.Get(u)
	doc, _ := goquery.NewDocumentFromReader(r.Body)
	defer r.Body.Close()
	h := doc.Find(".sample-tests .input")
	converter := md.NewConverter(S.Host, true, nil)
	reg := regexp.MustCompile("(?s)```(.*?)```")
	h.Each(func(i int, s *goquery.Selection) {
		res := reg.FindSubmatch([]byte(converter.Convert(s)))
		if len(res) > 0 {
			p.input = append(p.input, strings.TrimSpace(string(res[1])))
		}
	})
	h = doc.Find(".output")
	h.Each(func(i int, s *goquery.Selection) {
		res := reg.FindSubmatch([]byte(converter.Convert(s)))
		if len(res) > 0 {
			p.output = append(p.output, strings.TrimSpace(string(res[1])))
		}
	})
	statement := converter.Convert(doc.Find(".problem-statement"))
	p.statement = statement
	return p
}
