package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
)

func main() {
	// Set color red
	red := color.New(color.FgHiRed).SprintFunc()

	// Check number of arguments
	if len(os.Args) < 2 {
		fmt.Printf("\n%s gen CLI tool currently only supports the '%s' command for token generation. \n\n", red("Error:"), red("token"))
		os.Exit(1)
	}

	// Retrieve first argument
	command := strings.ToLower(os.Args[1])

	// Check if command is supported
	switch command {
	case "token":
		// Define and parse the -b flag
		tokenCmd := flag.NewFlagSet("token", flag.ExitOnError)
		bPtr := tokenCmd.Int("b", 0, "Number of bytes to generate")
		err := tokenCmd.Parse(os.Args[2:])
		if err != nil {
			fmt.Printf("\n%s Unable to generate random bytes - %s.\n\n", red("Error:"), err)
			return
		}

		// Generate token
		GenerateToken(*bPtr)
	default:
		// Print error message
		fmt.Printf("\n%s gen CLI tool currently only supports the '%s' command for token generation. \n\n", red("Error:"), red("token"))
		os.Exit(1)
	}
}
