package cmd

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

func Test() {
	path, err := os.Getwd()
	if err != nil {
		color.Red("Error getting %v\n%v", path, err.Error())
		return
	}
	files, err := os.ReadDir(path)
	if err != nil {
		color.Red("Error reading files in %v\n%v", path, err.Error())
		return
	}
	for _, file := range files {
		fmt.Println(file.Name())
	}
}
