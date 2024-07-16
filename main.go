package main

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hunderaweke/codative-codeforces/client"
	"golang.org/x/term"
)

func main() {
	client.Create("https://codeforces.com")
	var handle string
	color.Blue("%s", "Insert Your handle")
	fmt.Fscanln(os.Stdin, &handle)
	color.Blue("%s", "Insert Your password")
	password, err := term.ReadPassword(0)
	err = client.Clnt.Login(handle, string(password))
	if err != nil {
		color.Red("%v", err.Error())
		return
	}
	color.Green("Login Successful")
}
