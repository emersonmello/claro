package cmd

import (
	"github.com/emersonmello/claro/utils"
	"github.com/spf13/cobra"
)

// messageCmd represents the message command
var messageCmd = &cobra.Command{
	Use:     "message <\"commit message\">",
	Example: "claro config message \"Graded project, the file containing the grade is in the root directory\"",
	Short:   "define the commit message for grading",
	Args:    cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		utils.WriteConfigFile("commit_message", args[0], "Commit message has been set successfully!")
	},
}

func init() {
	configCmd.AddCommand(messageCmd)
}
