package cmd

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/emersonmello/claro/utils"
	"github.com/pterm/pterm"

	"github.com/spf13/cobra"
)

func getRepositoriesAndGradeFiles(files []fs.FileInfo) map[string]fs.FileInfo {
	fileNamePattern := "^(grade-).*(\\.md)$"
	regexpPattern, _ := regexp.Compile(fileNamePattern)

	oneFilePerRepo := make(map[string]fs.FileInfo)
	reposMissing := make([]string, 0)
	gradeMissing := make([]string, 0)

	for _, f := range files {
		if f.IsDir() {
			oneFilePerRepo[f.Name()] = nil
		}
	}

	for _, f := range files {
		if f.Mode().IsRegular() && regexpPattern.MatchString(f.Name()) {
			if _, present := oneFilePerRepo[strings.Split(strings.Split(f.Name(), "grade-")[1], ".md")[0]]; present {
				oneFilePerRepo[strings.Split(strings.Split(f.Name(), "grade-")[1], ".md")[0]] = f
			} else {
				reposMissing = append(reposMissing, strings.Split(strings.Split(f.Name(), "grade-")[1], ".md")[0])
			}
		}
	}

	for k, v := range oneFilePerRepo {
		if v == nil {
			gradeMissing = append(gradeMissing, k)
		}
	}

	if len(gradeMissing) > 0 {
		pterm.Warning.Println("Attention! There is no grade file for repositories below")
		for _, v := range gradeMissing {
			fmt.Println(" --> " + v)
		}
	}

	if len(reposMissing) > 0 {
		pterm.Warning.Println("Attention! There is no repositories for grade files below")
		for _, v := range reposMissing {
			fmt.Println(" --> " + v)
		}
	}

	return oneFilePerRepo
}

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push <assignment repositories prefix>",
	Short: "Add and commit the grading file in each student repository",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		files, err := ioutil.ReadDir(args[0])
		if err != nil {
			checkError(err)
		}

		oneFilePerRepo := getRepositoriesAndGradeFiles(files)

		e := os.Chdir(args[0])
		checkError(e)

		p := pterm.DefaultProgressbar.WithRemoveWhenDone(true)
		p.ShowPercentage = false
		p.ShowCount = false
		p.ShowCount = false
		p.ShowElapsedTime = false
		p.Start()
		for repo, gradeFile := range oneFilePerRepo {

			if gradeFile != nil {
				utils.AddAndCommitGradeFile(gradeFile.Name(), repo, p)
			}
		}
		p.Stop()
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)
}
