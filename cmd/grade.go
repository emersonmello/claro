/*
Copyright © 2022  Emerson Mello

*/
package cmd

import (
	"github.com/emersonmello/claro/utils"
	"github.com/spf13/cobra"
)

// gradeCmd represents the grade command
var gradeCmd = &cobra.Command{
	Use:     "grade <\"string\">",
	Example: "claro config grade \"Grade: \"",
	Short:   "define the grade string inserted in the file representing the grade sheet",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utils.WriteConfigFile("grade_string", args[0], "Grade string has been set successfully!")
	},
}

func init() {
	configCmd.AddCommand(gradeCmd)
}
