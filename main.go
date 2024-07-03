package main

import (
	"github.com/fatih/color"
	"github.com/hunderaweke/codative-codeforces/utils"
)

func main() {
	data, err := utils.Fetch()
	if err != nil {
		color.Red(err.Error())
		return
	}
	utils.CreateFiles(*data)
}
