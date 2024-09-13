// Package internal
package internal

/*
Copyright © 2022-2024 Emerson Ribeiro de Mello <mello@ifsc.edu.br>
*/

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/huh"
	"github.com/emersonmello/claro/internal/tui"
	"github.com/zalando/go-keyring"
)

const service = "a github classroom cli"
const user = "claro"

// DeletePasswordItem Delete the user's GitHub personal access token from the operating system keyring
func deletePasswordItem() error {
	return keyring.Delete(service, user)
}

// CreateKey Store the user's GitHub personal access token in the operating system keyring
func createKey(password string, removeIfExist bool) error {
	if removeIfExist {
		_ = deletePasswordItem()
	}
	e := keyring.Set(service, user, password)
	if e != nil {
		fmt.Println(tui.ErrorStyle.Render(fmt.Sprintf("Could not store token in operating system keyring:\n => %s", e)))
	} else {
		fmt.Println(tui.DoneStyle.Render("Your github personal access token has been successfully set in the operating system keyring!"))
	}
	return e
}

// GetPassword Retrieve the user's GitHub personal access token from the operating system keyring
func getPassword() (string, error) {
	return keyring.Get(service, user)
}

// ReadTokenFromStdIn To obtain the user's GitHub Personal Access Token
func readTokenFromStdIn() string {
	var userToken string
	group := huh.NewGroup(
		huh.NewInput().
			Value(&userToken).Placeholder("Ex: ghp_1873SsDhdjf....").
			Title("Provide your GitHub Personal Access Token (classic):").
			EchoMode(huh.EchoModePassword),
	)
	form := huh.NewForm(group)
	if err := form.Run(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
		os.Exit(0)
	}
	return userToken
}

// yesNoDialog displays a yes/no confirmation dialog with the given message.
// It returns true if the user confirms, otherwise false.
func yesNoDialog(msg string) bool {
	var confirm bool
	group := huh.NewGroup(
		huh.NewConfirm().
			Title(msg).
			Value(&confirm),
	)
	form := huh.NewForm(group)
	if e := form.Run(); e != nil {
		_, _ = fmt.Fprintln(os.Stderr, e)
		os.Exit(0)
	}
	return confirm
}

// DeleteTokenFromKeyring deletes the GitHub Personal Access Token from the OS keyring.
func DeleteTokenFromKeyring() {
	if ghToken, _ := getPassword(); ghToken != "" {
		confirm := yesNoDialog("Are you sure you want to delete the GitHub Personal Access Token from the operating system keyring?")
		if confirm {
			if err := deletePasswordItem(); err != nil {
				if strings.Contains(err.Error(), "secret not found") {
					fmt.Println(tui.ErrorStyle.Render("Secret not found in OS Keychain."))
				} else {
					fmt.Println(tui.ErrorStyle.Render(fmt.Sprintf("Error deleting token from OS Keychain: %s", err)))
				}
			} else {
				fmt.Println(tui.DoneStyle.Render("Token deleted from OS Keychain"))
			}
		}
	} else {
		fmt.Println(tui.ErrorStyle.Render("No token found in OS Keychain. Nothing to delete."))
	}
}

// AddTokenToKeyring adds a GitHub Personal Access Token to the OS keyring.
func AddTokenToKeyring() {
	if ghToken, _ := getPassword(); ghToken != "" {
		confirm := yesNoDialog("A GitHub Personal Access Token for claro is already stored in the operating system keyring. Would you like to override it?")
		if !confirm {
			return
		}
	}
	if ghToken := readTokenFromStdIn(); ghToken != "" {
		_ = createKey(ghToken, true)
	}
}

// GetAndSaveToken retrieves the GitHub Personal Access Token from the OS keyring.
// Returns the GitHub Personal Access Token.
func GetAndSaveToken() string {
	var ghToken string
	if ghToken, _ = getPassword(); ghToken == "" {
		if ghToken = readTokenFromStdIn(); ghToken != "" {
			// If the token is not found, it prompts the user to input the token and optionally saves it in the OS keyring.
			//persist := yesNoDialog("Would you like to save this token in the OS keyring?")
			//if persist {
			//	_ = createKey(ghToken, true)
			//}
		}
	}
	return ghToken
}

// writeConfigFile key/value in the config file
// func writeConfigFile(key, value, returnMessage string) {
// 	viper.Set(key, value)
// 	if viper.WriteConfig() != nil {
// 		err := viper.SafeWriteConfig()
// 		if err != nil {
// 			fmt.Println(tui.ErrorStyle.Render(fmt.Sprintf("Error writing config to file: %s", err)))
// 		}
// 	}
// 	fmt.Println(returnMessage)
// }
