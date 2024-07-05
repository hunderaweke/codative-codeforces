package utils

import (
	"fmt"
	"os"

	"github.com/fatih/color"
	"github.com/hunderaweke/codative-codeforces/types"
)

func CreateFiles(data types.Response) error {
	directoryName := reformString(data.Result.Contest.Name)

	err := os.Mkdir(string(directoryName), 0755)
	if err != nil {
		return err
	}
	err = os.Chdir(string(directoryName))
	if err != nil {
		return err
	}
	for _, prob := range data.Result.Problems {
		probName := reformString(prob.Index + " " + prob.Name)
		os.Mkdir(probName, 0755)
		os.Chdir(probName)
		color.Blue(fmt.Sprintf("Generating %s file and directories", prob.Index+" "+prob.Name))
		// TODO: Implement Template File Generating
		FetchTestCases(data.Result.Contest.Id, prob.Index, prob.Name)
		color.Green(fmt.Sprintf("Finished Generating %s file and directories", prob.Index+" "+prob.Name))
		os.Chdir("..")
	}
	color.Green(fmt.Sprintf("Finished Creating the directories and files for %s", data.Result.Contest.Name))
	return nil
}
