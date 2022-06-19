/*
Copyright © 2022  Emerson Mello

*/
package cmd

import (
	"github.com/emersonmello/claro/utils"

	"os"

	"github.com/spf13/cobra"
)

// pullCmd represents the pull command
var pullCmd = &cobra.Command{
	Use:   "pull <assignment repositories prefix>",
	Short: "Incorporate changes from students' remote repositories into local copy",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		files, err := os.ReadDir(args[0])
		if err != nil {
			checkError(err)
		}
		utils.Pull(files, args[0])
	},
}

func init() {
	rootCmd.AddCommand(pullCmd)
}
