package tools

import (
	"fmt"
	"strings"
)

func AskForConfirmation(question string) (bool, error) {
	var response string

	fmt.Print(question)
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false, err
	}

	switch strings.ToLower(response) {
	case "y", "yes":
		return true, nil
	case "n", "no":
		return false, nil
	default:
		fmt.Println("I'm sorry but I didn't get what you meant, please type (y)es or (n)o and then press enter:")
		return AskForConfirmation(question)
	}
}
