package cmd

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/hunderaweke/codative-codeforces/internal"
	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/shirou/gopsutil/process"
)

var Extensions map[string]bool

func init() {
	Extensions = make(map[string]bool)
	for _, ext := range internal.FileExtensions {
		Extensions[ext] = true
	}
}

func Test() error {
	path, err := os.Getwd()
	if err != nil {
		return err
	}
	entry, _ := os.ReadDir(path)
	fileName := ""
	inputs := make(map[int64]*os.File)
	outputs := make(map[int64]string)
	for _, e := range entry {
		if !e.IsDir() {
			ext := filepath.Ext(e.Name())
			if ext == ".py" {
				fileName = e.Name()
			}
			reg := regexp.MustCompile(`input(\d+?)\.in`)
			found := reg.FindSubmatch([]byte(e.Name()))
			if len(found) > 1 {
				num, _ := strconv.ParseInt(string(found[1]), 10, 64)
				file, _ := os.Open(string(found[0]))
				inputs[num] = file
				defer file.Close()
			}
			reg = regexp.MustCompile(`output(\d+?)\.out`)
			found = reg.FindSubmatch([]byte(e.Name()))
			if len(found) > 1 {
				num, _ := strconv.ParseInt(string(found[1]), 10, 64)
				file, _ := os.ReadFile(string(found[0]))
				outputs[num] = string(file)
			}
			/* _, ok := cmd.Extensions[ext]
			if ok {
				fmt.Println(e.Name())
			} */
		}
	}
	for id, input := range inputs {
		err = testSample(int(id+1), input, internal.Lang{Name: "python", Command: "python %v"}, fileName, outputs[id])
		if err != nil {
			return err
		}
	}
	return err
}

func testSample(sampleID int, sample io.Reader, lang internal.Lang, fileName string, expectedOutput string) error {
	cmd := exec.Command(lang.Name, fileName)
	var outputBuffer bytes.Buffer
	cmd.Stdin = sample
	cmd.Stdout = &outputBuffer
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("Runtime Error #%v ..... %v", sampleID, err.Error())
	}
	pid := int32(cmd.Process.Pid)
	maxMemory := uint64(0)
	ch := make(chan error)
	go func() {
		ch <- cmd.Wait()
	}()
	running := true
	for running {
		select {
		case err := <-ch:
			if err != nil {
				return fmt.Errorf("Runtime Error #%v .... %v", sampleID, err.Error())
			}
			running = false
		default:
			p, err := process.NewProcess(pid)
			if err == nil {
				m, err := p.MemoryInfo()
				if err == nil {
					if m.RSS > maxMemory {
						maxMemory = m.RSS
					}
				}
			}
		}
	}
	foundOutput := strings.TrimSuffix(outputBuffer.String(), "")
	similar := strings.Compare(foundOutput, expectedOutput)
	if similar == 0 {
		memory := ""
		if maxMemory >= 1024*1024 {
			memory = fmt.Sprintf("%.3f MB", float64(maxMemory)/(1024.0*1024.0))
		} else if maxMemory >= 1024 {
			memory = fmt.Sprintf("%.3f KB", float64(maxMemory)/1024.0)
		} else {
			memory = fmt.Sprintf("%+v Bytes", maxMemory)
		}
		color.Green("Passed #%v...%v  %vs", sampleID, memory, cmd.ProcessState.UserTime().Seconds())
		return nil
	}
	dmp := diffmatchpatch.New()
	diff := dmp.DiffMain(foundOutput, expectedOutput, true)
	red := color.New(color.FgRed).Add(color.Bold).PrintfFunc()
	redUnderline := color.New(color.FgHiRed).Add(color.Underline, color.Bold).PrintfFunc()
	blueUnderline := color.New(color.FgHiBlue).Add(color.Underline, color.Bold).PrintfFunc()
	green := color.New(color.FgGreen).Add(color.Bold).PrintfFunc()
	redUnderline("Failed on test #%v\n", sampleID)
	blueUnderline("Expected\tFound\n")
	idx := 0
	for idx < len(diff) {
		d := diff[idx]
		if d.Type.String() == "Delete" {
			wrong := d.Text
			idx++
			d = diff[idx]
			green("%v\t\t", d.Text)
			red("%v\n", wrong)
		} else {
			match := strings.TrimPrefix(strings.TrimSuffix(d.Text, "\n"), "\n")
			green("%v\t\t%v\n", match, match)
		}
		idx += 1
	}
	return nil
}
