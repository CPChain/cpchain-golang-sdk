package tools

import (
	"fmt"
	"os"

	"github.com/gookit/color"
	"github.com/howeyc/gopass"
	"github.com/zgljl2012/slog"
)

func GetPassword(prompt string, needConfirm bool) (string, error) {
	if prompt != "" {
		color.HiGreen.Println(prompt)
	}
	passwordByte, err := gopass.GetPasswdPrompt("Password: ", false, os.Stdin, os.Stdout)
	password := string(passwordByte)
	if err != nil {
		color.HiRed.Printf("ERROR: Failed to read password: %v", err)
		os.Exit(1)
	}
	if needConfirm {
		confirmByte, err := gopass.GetPasswdPrompt("Repeat password: ", true, os.Stdin, os.Stdout)
		confirm := string(confirmByte)
		if err != nil {
			color.HiRed.Printf("ERROR: Failed to read password confirmation: %v", err)
			os.Exit(1)
		}
		if password != confirm {
			color.HiRed.Printf("ERROR: Password do not match")
			os.Exit(1)
		}
	}
	return password, nil
}

func AskForConfirmation(prompt string) bool {
	if prompt != "" {
		color.HiBlue.Println(prompt)
	}
	return askForConfirmation()
}

func askForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		slog.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if containsString(okayResponses, response) {
		return true
	} else if containsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return askForConfirmation()
	}
}

func containsString(items []string, item string) bool {
	for _, eachItem := range items {
		if eachItem == item {
			return true
		}
	}
	return false
}
