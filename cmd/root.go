/*
Copyright © 2022  Emerson Mello

*/
package cmd

import (
	"os"

	"github.com/emersonmello/claro/utils"
	"github.com/pterm/pterm"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var (
	rootCmd = &cobra.Command{
		Use:     "claro",
		Short:   "A GitHub Classroom CLI for teachers",
		Version: "0.1.0",
	}
)

func checkError(e error) {
	if e != nil {
		pterm.Error.Println(e.Error())
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	utils.CheckExternalsCommands()
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.CompletionOptions.DisableDescriptions = true
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config-file", "", "config file (default is $HOME/.claro.env)")
}

func initConfig() {
	utils.GetConf()
}
