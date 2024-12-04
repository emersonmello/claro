// Package cmd
package cmd

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime/debug"

	"github.com/emersonmello/claro/cmd/clone"
	"github.com/emersonmello/claro/cmd/config"
	"github.com/emersonmello/claro/cmd/pull"
	"github.com/emersonmello/claro/cmd/push"
	"github.com/emersonmello/claro/cmd/token"
	"github.com/emersonmello/claro/internal"
	"github.com/emersonmello/claro/internal/tui"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var pathConfigFile string

var version = "1.0.1"

// Print program version
func programVersion() string {
	result := version
	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				commit := setting.Value[0:7]
				result = fmt.Sprintf("%s\ncommit: %s", result, commit)
			}
			if setting.Key == "vcs.time" {
				date := setting.Value[0:10]
				result = fmt.Sprintf("%s, built at: %s", result, date)
			}
		}
	}
	return result
}

var rootCmd = &cobra.Command{
	Use:     "claro",
	Short:   "A GitHub Classroom CLI for teachers",
	Version: programVersion(),
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	str := fmt.Sprintf("config file (default is %s/config.env)", internal.ConfigDir())

	rootCmd.PersistentFlags().StringVar(&cfgFile,
		"config",
		"",
		str)

	// Add subcommands

	pullCmd := pull.Pull()
	pullCmd.Example = fmt.Sprintf("%s %s assignment-01-submissions", rootCmd.CommandPath(), pullCmd.Name())

	pushCmd := push.Push()
	pushCmd.Example = fmt.Sprintf("%s %s assignment-01-submissions", rootCmd.CommandPath(), pushCmd.Name())

	tokenCmd := token.Token()
	tokenCmd.Example = fmt.Sprintf("%s %s add\n%s %s del", rootCmd.CommandPath(), tokenCmd.Name(), rootCmd.CommandPath(), tokenCmd.Name())

	rootCmd.AddCommand(clone.Clone())
	rootCmd.AddCommand(config.Config())
	rootCmd.AddCommand(tokenCmd)
	rootCmd.AddCommand(pullCmd)
	rootCmd.AddCommand(pushCmd)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory $HOME/.config/claro/config.env
		pathConfigFile = internal.ConfigDir()
		viper.AddConfigPath(pathConfigFile)
		viper.SetConfigType("env")
		viper.SetConfigName("config")
	}

	viper.SetDefault("version", internal.ClaroConfigStrings.Version)
	viper.SetDefault("message", internal.ClaroConfigStrings.Message)
	viper.SetDefault("filename", internal.ClaroConfigStrings.Filename)
	viper.SetDefault("title", internal.ClaroConfigStrings.Title)
	viper.SetDefault("grade", internal.ClaroConfigStrings.Grade)

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		var configFileNotFoundError viper.ConfigFileNotFoundError
		// Config file not found
		if errors.As(err, &configFileNotFoundError) {
			// Creating config directory
			if _, err = os.Stat(pathConfigFile); os.IsNotExist(err) {
				if e := os.MkdirAll(pathConfigFile, 0755); e != nil {
					_, _ = fmt.Fprintln(os.Stderr, err)
				}
			}
		}
		// Creating config file
		if err = viper.SafeWriteConfig(); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
		}
	}

	// unmarshal config and storing it on runtime conf var
	if err := viper.Unmarshal(internal.ClaroConfigStrings); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}

	if !checkGitAndCredentials() {
		os.Exit(1)
	}
}

func checkGitAndCredentials() bool {
	if _, err := exec.LookPath("git"); err != nil {
		fmt.Println(tui.ErrorStyle.Render("I can't find 'git' command. Please, be sure that 'git' is installed and in the user PATH"))
		return false
	}
	dirname, _ := os.UserHomeDir()
	cmd := exec.Command("git", "config", "--global", "--get", "credential.helper")
	cmd.Dir = dirname
	if out, _ := cmd.Output(); out != nil {
		if len(out) == 0 {
			cmd = exec.Command("git", "config", "--global", "--add", "credential.helper", "cache")
			cmd.Dir = dirname
			_ = cmd.Run()
		}
	}
	// Checking if you have GitHub CLI installed
	if _, err := exec.LookPath("gh"); err == nil {
		tui.GitHubCliInstalled = true
	}
	return true
}
