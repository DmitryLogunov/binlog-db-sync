package errors

import (
	"fmt"
	"os"
)

func CheckAndExitIfError(err error, info string) {
	if err != nil {
		fmt.Printf("\nError: \n")
		fmt.Println(err)
		fmt.Printf("\nInfo: %s\n\n", info)
		os.Exit(1)
	}
}

func CheckAndReturnIfError(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}

	return false
}
