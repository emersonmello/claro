package cmd

import (
	"fmt"

	"github.com/emersonmello/claro/utils"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:           "list",
	Short:         "List all student assignment repositories in an organization",
	Example:       "claro list ifsc-classroom 2022-01-assignment-01",
	SilenceErrors: true,
	Args:          cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		repositories := utils.GetRepositoryList(args[0], args[1])

		if len(repositories) > 0 {
			fmt.Printf("%d repositories were found with %s prefix!\n", len(repositories), args[1])
			for _, r := range repositories {
				fmt.Println(r.Name)
			}
		} else {
			fmt.Printf("Organization %s does not contains repositories with prefix %s\n", args[0], args[1])
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
