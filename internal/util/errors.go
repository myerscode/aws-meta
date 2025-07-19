package util

import (
	"fmt"
	"os"
)

func PrintErrorAndExit(err error) {
	fmt.Printf("error: %v\n", err)
	os.Exit(1)
}
