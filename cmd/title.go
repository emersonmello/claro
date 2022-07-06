package cmd

import (
	"github.com/emersonmello/claro/utils"
	"github.com/spf13/cobra"
)

// titleCmd represents the title command
var titleCmd = &cobra.Command{
	Use:     "title <\"string\">",
	Example: "claro config title \"Feedback\"",
	Short:   "define the title string in the file representing the grade sheet",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utils.WriteConfigFile("grade_title", args[0], "Title string has been set successfully!")
	},
}

func init() {
	configCmd.AddCommand(titleCmd)
}
