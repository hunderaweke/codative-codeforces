package internal

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/hunderaweke/codative-codeforces/session"
	"github.com/hunderaweke/codative-codeforces/utils"
)

type problem struct {
	title  string
	input  []string
	output []string
}

func (p *problem) create(data []byte, ext string) error {
	fmt.Printf("Generating %s\n", p.title)
	problemDir := utils.ReformString(p.title)
	os.Mkdir(problemDir, 0777)
	os.Chdir(problemDir)
	for i, input := range p.input {
		name := fmt.Sprintf("input%d.in", i)
		file, err := os.Create(name)
		if err != nil {
			return err
		}
		defer file.Close()
		file.Write([]byte(input))
	}
	for i, output := range p.output {
		name := fmt.Sprintf("output%d.out", i)
		file, err := os.Create(name)
		if err != nil {
			return err
		}
		defer file.Close()
		file.Write([]byte(output))
	}
	file, err := os.Create(problemDir + ext)
	defer file.Write(data)
	if err != nil {
		return err
	}
	os.Chdir("../")
	fmt.Printf("Done %s\n", p.title)
	return nil
}

func parseProblem(contestID, contestType, problemID string) (p problem) {
	u := strings.Join([]string{session.S.Host, contestType, contestID, "problem", problemID}, "/")
	r, _ := session.S.Client.Get(u)
	doc, _ := goquery.NewDocumentFromReader(r.Body)
	defer r.Body.Close()
	h := doc.Find(".sample-tests .input")
	converter := md.NewConverter(session.S.Host, true, nil)
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
	return p
}
