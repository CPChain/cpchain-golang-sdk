package tools

import (
	"os"

	"github.com/gookit/color"
	"github.com/howeyc/gopass"
)

func GetPassword(prompt string, needConfirm bool) (string, error) {
	if prompt != "" {
		color.HiGreen.Println(prompt)
	}
	passwordByte, err := gopass.GetPasswdPrompt("Password: ", true, os.Stdin, os.Stdout)
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

func NeedConfirm() {

}
