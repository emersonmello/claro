// Package internal
package internal

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	xdgConfigHome = "XDG_CONFIG_HOME"
)

type ClaroCfg struct {
	Version  int    `mapstructure:"version"`
	Message  string `mapstructure:"message"`
	Filename string `mapstructure:"filename"`
	Title    string `mapstructure:"title"`
	Grade    string `mapstructure:"grade"`
}
type choice int

const (
	filename choice = iota
	message
	title
	grade
	quit
)

const configFilename = "claro"

func ConfigDir() string {
	var path string

	if l := os.Getenv(xdgConfigHome); l != "" {
		path = filepath.Join(l, configFilename)
	} else if a := os.Getenv("AppData"); runtime.GOOS == "windows" && a != "" {
		path = filepath.Join(a, configFilename)
	} else {
		home, _ := os.UserHomeDir()
		path = filepath.Join(home, ".config", configFilename)
	}
	return path
}

var ClaroConfigStrings = &ClaroCfg{
	Version:  1,
	Message:  "This project has been graded. The file containing the grade is located in the root directory.",
	Filename: "GRADING.md",
	Title:    "Feedback",
	Grade:    "Grade: ",
}

func ConfigCmd(cmd *cobra.Command, args []string) error {

	var option choice

	for {

		form := huh.NewForm(
			huh.NewGroup(
				huh.NewSelect[choice]().
					Title("Which option do you want to configure?").
					Options(
						huh.NewOption("Grade filename", filename),
						huh.NewOption("Commit message", message),
						huh.NewOption("Grade file's title", title),
						huh.NewOption("Grade file's grade string", grade),
						huh.NewOption("Quit", quit),
					).
					Value(&option),
			),
		)

		if err := form.Run(); err != nil {
			fmt.Println("There was an error running the program:", err)
		}

		var group *huh.Group

		switch option {
		case filename:
			group = huh.NewGroup(
				huh.NewInput().
					Value(&ClaroConfigStrings.Filename).
					Title("The name of the file that will be created in the student repository containing the feedback."),
			)
		case message:
			group = huh.NewGroup(
				huh.NewInput().
					Value(&ClaroConfigStrings.Message).
					Title("Commit message for grading"),
			)
		case title:
			group = huh.NewGroup(
				huh.NewInput().
					Value(&ClaroConfigStrings.Title).
					Title("The file's title representing the grade sheet"),
			)
		case grade:
			group = huh.NewGroup(
				huh.NewInput().
					Value(&ClaroConfigStrings.Grade).
					Title("The grade string inserted in the file representing the grade sheet."),
			)
		case quit:
			// Saving config file
			viper.Set("Title", ClaroConfigStrings.Title)
			viper.Set("Message", ClaroConfigStrings.Message)
			viper.Set("Filename", ClaroConfigStrings.Filename)
			viper.Set("Grade", ClaroConfigStrings.Grade)
			if err := viper.WriteConfig(); err != nil {
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
			return nil
		}

		form = huh.NewForm(group)

		if err := form.Run(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
	}
}
