package cmd

import (
	"github.com/emersonmello/claro/utils"
	"github.com/spf13/cobra"
)

// filenameCmd represents the filename command
var filenameCmd = &cobra.Command{
	Use:     "filename <filename.md>",
	Example: "claro config filename \"GRADING.md\" ",
	Short:   "define the name of the file that will be created in the student repository containing the feedback",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utils.WriteConfigFile("grade_filename", args[0], "Grade sheet filename has been set successfully!")
	},
}

func init() {
	configCmd.AddCommand(filenameCmd)
}
