package cmd

import (
	"fmt"

	"github.com/emersonmello/claro/utils"

	"github.com/spf13/cobra"
)

// tokenCmd represents the token command
var tokenCmd = &cobra.Command{
	Use:   "token <add|del>",
	Short: "Add or remove GitHub Personal Access Token in OS Keychain",
	Long: `Add or remove GitHub Personal Access Token in OS Keychain

If the OS Keychain is not available then you can store the token in the config file`,
	Example:   "claro config token add",
	ValidArgs: []string{"add", "del"},
	Args:      cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
	Run: func(cmd *cobra.Command, args []string) {
		switch args[0] {
		case "add":
			utils.SaveGHToken(utils.ReadTokenFromStdIn())
		case "del":
			if !utils.PromptUser("Do you want to remove the token from the keychain? (y/N):", "no") {
				if e := utils.DeletePasswordItem(); e == nil {
					fmt.Println("Done!")
				} else {
					checkError(e)
				}
			}
		}
	},
}

func init() {
	configCmd.AddCommand(tokenCmd)
}
