// A set of utility functions to handle with git command
package utils

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"

	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

// To clone repositories from GitHub
func CloneRepositories(outputDir string, repositories []RepData) ([]string, []string) {

	successCloned := make([]string, 0)
	errorCloned := make([]string, 0)

	p, _ := pterm.DefaultProgressbar.WithTotal(len(repositories)).WithRemoveWhenDone(true).WithTitle("Cloning students repositories").Start()

	for _, r := range repositories {

		p.UpdateTitle("Cloning " + r.Name + " ")
		cmd := exec.Command("git", "clone", r.Url)
		cmd.Dir = outputDir
		if err := cmd.Run(); err == nil {
			pterm.Success.Println("Cloning " + r.Name)
			successCloned = append(successCloned, r.Name)
		} else {
			pterm.Error.Printfln("To clone " + r.Name)
			errorCloned = append(errorCloned, err.Error())
		}
		p.Increment()
	}
	return successCloned, errorCloned
}

// To add and commit a grade file in a GitHub repository
func AddAndCommitGradeFile(filename, repositoryDir string, p *pterm.ProgressbarPrinter) {
	gradeFileName := viper.GetString("grade_filename")

	p.UpdateTitle("Push '" + gradeFileName + "' to " + repositoryDir)
	if src, err := os.Open(filename); err == nil {
		defer src.Close()
		if err := os.Chdir(repositoryDir); err == nil {
			if dst, err := os.Create(gradeFileName); err == nil {
				defer dst.Close()
				if _, err = io.Copy(dst, src); err == nil {
					execGitCommands("reset")
					execGitCommands("add", gradeFileName)
					if o, err := execGitCommands("commit", "-m", viper.GetString("commit_message")); err == nil {
						execGitCommands("pull", "--rebase")
						execGitCommands("push")
						pterm.Success.Println("Push '" + gradeFileName + "' to " + repositoryDir)
					} else {
						pterm.Warning.Println(repositoryDir + " has a problem!")
						fmt.Println(string(o))
					}
				} else {
					pterm.Error.Println("Could not copy the file: " + filename + " to " + repositoryDir + "/" + gradeFileName)
				}
			} else {
				pterm.Error.Println("Could not create the file: " + filename)
			}
			os.Chdir("..")
		} else {
			pterm.Error.Println("Could not access the directory: " + repositoryDir)
		}
	} else {
		pterm.Error.Println("Could not open the file: " + filename)
	}
	p.UpdateTitle("")
}

// To pull all GitHub repositories inside a specific local directory
func Pull(files []fs.DirEntry, directory string) {
	e := os.Chdir(directory)
	checkError(e)
	p := pterm.DefaultProgressbar.WithRemoveWhenDone(true)
	p.ShowPercentage = false
	p.ShowCount = false
	p.ShowCount = false
	p.ShowElapsedTime = false
	p.Start()
	for _, f := range files {
		if f.IsDir() {
			p.UpdateTitle("Pull " + f.Name())
			execGitCommands("stash")
			cmd := exec.Command("git", "pull", "--rebase")
			cmd.Dir = f.Name()
			_, err := cmd.CombinedOutput()
			if err != nil {
				pterm.Error.Println(f.Name() + " is not clean!")
				pterm.Info.Println(e.Error())
			} else {
				execGitCommands("stash", "pop")
				pterm.Success.Println("Pull " + f.Name())
			}
		}
	}
	p.Stop()
}

// A wrapper to execute an external command - git
func execGitCommands(pars ...string) ([]byte, error) {
	cmd := exec.Command("git", pars...)
	return cmd.CombinedOutput()
}

// check error
func checkError(e error) {
	if e != nil {
		fmt.Printf("Error: %s\n", e)
		os.Exit(1)
	}
}
