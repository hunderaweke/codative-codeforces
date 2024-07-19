package internal

type problem struct {
	title     string
	statement string
	input     []string
	output    []string
}
func (p *problem) create(tmplId int) error {
	problemDir := utils.ReformString(p.title)
	os.Mkdir(problemDir, 0644)
	os.Chdir(problemDir)
	for i, input := range p.input {
		file, err := os.Create(string(i) + ".in")
		if err != nil {
			return err
		}
		file.Write([]byte(input))
		file.Close()
	}
	for i, output := range p.output {
		file, err := os.Create(string(i) + ".in")
		if err != nil {
			return err
		}
		file.Write([]byte(output))
		file.Close()
	}
	template := 
	return nil
}

func parseProblem(contestID, contestType, problemID string) (p problem) {
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
