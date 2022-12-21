// A set of utility functions to handle with git command
package utils

import (
	"bytes"
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

					exec.Command("git", "reset").Run()
					exec.Command("git", "add", gradeFileName).Run()

					if o, _ := exec.Command("git", "status", "--porcelain").CombinedOutput(); string(o) == "" {
						pterm.Warning.Println(repositoryDir + ": nothing to commit, working tree clean")
					} else {

						if e := exec.Command("git", "commit", "-m", viper.GetString("commit_message")).Run(); e == nil {
							if e := exec.Command("git", "pull", "--rebase").Run(); e != nil {
								pterm.Error.Println(repositoryDir + ": could not execute git pull")
							} else {
								if e := exec.Command("git", "push", "--porcelain").Run(); e != nil {
									pterm.Error.Println(repositoryDir + ": could not execute git push")
								} else {
									pterm.Success.Println("Push '" + gradeFileName + "' to " + repositoryDir)
								}
							}
						} else {
							pterm.Error.Println(repositoryDir + ": could not execute git commit")
						}
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

	previousPrefix := pterm.Success.Prefix
	newCommitPrefix := pterm.Prefix{Text: "NEW COMMITS", Style: pterm.NewStyle(pterm.BgLightYellow, pterm.FgBlack)}

	for _, f := range files {
		if f.IsDir() {
			p.UpdateTitle("Pull " + f.Name())
			cmd := exec.Command("git", "stash", "--include-untracked")
			cmd.Dir = f.Name()
			err := cmd.Run()
			if err != nil {
				pterm.Error.Println("Could not run git stash command. Is the '" + f.Name() + "' directory a git repository?")
			} else {
				cmd = exec.Command("git", "rev-list", "--all", "--count")
				cmd.Dir = f.Name()
				totalCommitsBeforePull, _ := cmd.Output()
				cmd = exec.Command("git", "pull", "--rebase")
				cmd.Dir = f.Name()
				err = cmd.Run()
				if err != nil {
					pterm.Error.Println(f.Name() + " is not clean!")
				} else {
					cmd = exec.Command("git", "rev-list", "--all", "--count")
					cmd.Dir = f.Name()
					totalCommitsAfterPull, _ := cmd.Output()
					cmd := exec.Command("git", "stash", "pop")
					cmd.Dir = f.Name()
					_ = cmd.Run()
					if !bytes.Equal(totalCommitsAfterPull, totalCommitsBeforePull) {
						pterm.Success.Prefix = newCommitPrefix
					}
					pterm.Success.Println("Pull " + f.Name())
					pterm.Success.Prefix = previousPrefix
				}
			}
		}
	}
	fmt.Println()
}

// check error
func checkError(e error) {
	if e != nil {
		fmt.Printf("Error: %s\n", e)
		os.Exit(1)
	}
}
