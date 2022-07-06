package cmd

import (
	"fmt"
	"os"

	"github.com/emersonmello/claro/utils"
	"github.com/pterm/pterm"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var OutputDir string

var cloneCmd = &cobra.Command{
	Use:           "clone <organization> <assignment repositories prefix>",
	Short:         "Clone all students assignment repositories in an organization",
	Example:       "claro clone ifsc-classroom 2022-01-assignment-01",
	SilenceErrors: true,
	Args:          cobra.ExactArgs(2),

	Run: func(cmd *cobra.Command, args []string) {

		if OutputDir == "repository prefix" {
			OutputDir = args[1]
		}

		if _, err := os.Stat(OutputDir); !os.IsNotExist(err) {
			pterm.Error.Printf("The output directory '%s' already exists on the current directory!\n", OutputDir)
			os.Exit(1)
		}

		repositories := utils.GetRepositoryList(args[0], args[1])

		if len(repositories) <= 0 {
			pterm.Warning.Printfln("Organization %s does not contains repositories with prefix %s\n", args[0], args[1])
			os.Exit(0)
		} else {
			fmt.Printf("%d repositories were found\n", len(repositories))
		}

		if !utils.PromptUser("Would you like to proceed (Y/n)? ", "yes") {
			os.Exit(0)
		}

		s, _ := pterm.DefaultSpinner.Start("Create " + OutputDir + " directory")
		if e := os.Mkdir(OutputDir, os.ModePerm); e != nil {
			OutputDir = ""
			s.Fail(e.Error())
			os.Exit(1)
		}
		s.Success()

		utils.CloneRepositories(OutputDir, repositories)

		e := os.Chdir(OutputDir)
		checkError(e)

		s, _ = pterm.DefaultSpinner.Start("Create Markdown files")
		for _, r := range repositories {
			f, e := os.Create("grade-" + r.Name + ".md")
			if e != nil {
				s.Fail(e.Error())
			}
			_, e = f.WriteString("# " + viper.GetString("grade_title") + "\n\n- ...\n- " + viper.GetString("grade_string") + " \n\n")
			defer f.Close()
			if e != nil {
				s.Fail(e.Error())
			}
		}
		s.Success()
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().StringVarP(&OutputDir, "output-directory", "o", "repository prefix", "Directory where the cloned repositories should be stored")
}
