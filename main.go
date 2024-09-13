/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>

claro is a GitHub Classroom CLI tool that provides a simple interface for teachers to clone all student repositories at once for grading and then send grades to all these repositories simultaneously.

Usage:

	claro [command]

Available Commands:

	clone       Clone all students assignments from a GitHub Classroom
	completion  Generate the autocompletion script for the specified shell
	config      Configure claro's properties (commit message, filename, etc)
	help        Help about any command
	pull        Incorporate changes from students' remote repositories into local copy
	push        Add, commit, and push the grading file to each student's remote repository
	token       add or remove a claro's GitHub Personal Access Token in the OS Keychain
*/
package main

import (
	"github.com/emersonmello/claro/cmd"
)

func main() {
	_ = cmd.Execute()
}
