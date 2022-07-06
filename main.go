/*
Copyright © 2022  Emerson Ribeiro de Mello

claro is a GitHub Classroom CLI tool that offers a simple interface that allows the teacher to clone all student repositories at once for grading and then send grades at once to all these repositories.

Usage:

claro [command]

Available Commands:
  clone       Clone all students assignment repositories in an organization
  config      Configure claro's properties (github token, commit message, etc.)
  help        Help about any command
  list        List all student assignment repositories in an organization
  pull        Incorporate changes from students' remote repositories into local copy
  push        Add and commit the grading file in each student repository

claro consumes Github's REST API to fetch the list of repositories (assignments) from a GitHub Classroom organization. So, claro uses a GitHub Personal Access Token.

claro will try to get GitHub Personal Access Token from: (1) operating system keyring; (2) environment var (GH_TOKEN); (3) claro's config file (default $HOME/.claro.env)
*/
package main

import (
	"github.com/emersonmello/claro/cmd"
)

func main() {
	cmd.Execute()
}

// func init() {
// 	c := make(chan os.Signal)
// 	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
// 	go func() {
// 		// handling CTRL + C signal
// 		<-c
// should I remove the output directory created by clone command?
// if yes, use ioutil.ReadDir() and os.RemoveAll()
// 		os.Exit(1)
// 	}()
// }
