package utils

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/viper"
)

// Struct to store config file vars
type Config struct {
	CommitMessage string `mapstructure:"COMMIT_MESSAGE"`
	GradeFileName string `mapstructure:"GRADE_FILENAME"`
	GradeTitle    string `mapstructure:"GRADE_TITLE"`
	GradeString   string `mapstructure:"GRADE_STRING"`
	GHToken       string `mapstructure:"GH_TOKEN"`
}

func CheckExternalsCommands() {
	if _, err := exec.LookPath("git"); err != nil {
		pterm.Error.Println("I can't find 'git' command. Please, be sure that 'git' is installed and in the user PATH")
		os.Exit(1)
	}
}

// A Yes/No prompt for user
func PromptUser(message string, question string) bool {
	regexpYes := "([Yy](es)?|^$)"
	regexpNo := "([Nn](o)?|^$)"
	var match bool

	reader := bufio.NewReader(os.Stdin)
	fmt.Print(message)
	text, _ := reader.ReadString('\n')
	text = strings.Replace(text, "\n", "", -1)

	switch question {
	case "yes":
		match, _ = regexp.MatchString(regexpYes, text)
	case "no":
		match, _ = regexp.MatchString(regexpNo, text)
	}
	return match
}

// Write key/value in the config file
func WriteConfigFile(key, value, returnMessage string) {
	viper.Set(key, value)
	if viper.WriteConfig() != nil {
		viper.SafeWriteConfig()
	}
	fmt.Println(returnMessage)
}

// user viper to read envvar and returns a Config struct
func GetConf() *Config {
	viper.SetConfigType("env")
	home, _ := os.UserHomeDir()
	viper.AddConfigPath(home)
	viper.SetConfigName(".claro")

	conf := Config{}

	viper.SetDefault("grade_filename", "GRADING.md")
	viper.SetDefault("grade_title", "Feedback")
	viper.SetDefault("grade_string", "Grade: ")
	viper.SetDefault("commit_message", "Graded project, the file containing the grade is in the root directory")

	// read in environment variables that match
	viper.AutomaticEnv()
	viper.ReadInConfig()

	viper.Unmarshal(&conf)
	return &conf
}
