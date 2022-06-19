package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/viper"

	"github.com/zalando/go-keyring"
)

const service = "a github classroom cli"
const user = "claro"

// Delete password from os keyring
func DeletePasswordItem() error {
	return keyring.Delete(service, user)
}

// Store a password in os keyring
func CreateKey(password string, removeIfExist bool) error {
	if removeIfExist {
		DeletePasswordItem()
	}
	return keyring.Set(service, user, password)
}

// Get the password from os keyring
func GetPassword() (string, error) {
	secret, err := keyring.Get(service, user)
	return secret, err
}

// trying to get the github personal access token
func GetAndSaveToken(save bool) string {

	// form operating system keyring
	ghToken, err := GetPassword()

	if ghToken == "" || err != nil {
		// from environment var
		ghToken = viper.GetString("gh_token")
		if ghToken == "" {
			// ok, user should provides it right now
			ghToken = ReadTokenFromStdIn()
			if !save {
				if !PromptUser("Would you like to save this token in the os keyring or in the config file?(Y/n)? ", "yes") {
					return ghToken
				}
			}
			SaveGHToken(ghToken)
		} else {
			fmt.Println("Got GitHub Personal Access Token from envvar or claro config file")
		}
	} else {
		fmt.Println("Got GitHub Personal Access Token from OS keyring")
	}
	return ghToken
}

// Ask user about gh token
func ReadTokenFromStdIn() string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Inform your github personal access token: ")
	text, _ := reader.ReadString('\n')
	// convert CRLF to LF
	return strings.Replace(text, "\n", "", -1)
}

// try to save the gh token: (1) os keyring; (2) claro config file
func SaveGHToken(ghToken string) {
	tokenKC, err := GetPassword()
	if tokenKC != "" {
		if PromptUser("A github personal access token for claro is already stored in the operating system keyring.\nWould you like to override it? (y/N): ", "no") {
			return
		}
	}
	if err != nil {
		e := CreateKey(ghToken, true)
		if e == nil {
			fmt.Println("Your github personal access token has been successfully set in the operating system keyring!")
		} else {
			if PromptUser("Could not store token in operating system keyring.\nWould you like to save it in the config file?(Y/n):", "yes") {
				WriteConfigFile("gh_token", ghToken, "Your github personal access token has been successfully set in the configuration file!")
			}
		}
	} else {
		if PromptUser("Could not store token in operating system keyring.\nWould you like to save it in the config file?(Y/n):", "yes") {
			WriteConfigFile("gh_token", ghToken, "Your github personal access token has been successfully set in the configuration file!")
		}
	}
}
