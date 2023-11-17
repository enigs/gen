package main

import (
	"embed"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
)

// Define global color functions
var (
	red    = color.New(color.FgHiRed).SprintFunc()
	yellow = color.New(color.FgHiYellow).SprintFunc()
	blue   = color.New(color.FgHiBlue).SprintFunc()
	cyan   = color.New(color.FgHiCyan).SprintFunc()
)

//go:embed framework.20231117.zip
var content embed.FS
var fileName = "framework.20231117.zip"

func main() {
	// Check number of arguments
	if len(os.Args) < 2 {
		format := "\n%s gen CLI tool currently only supports '%s' command for token generation or '%s' for framework generation. \n\n"
		fmt.Printf(format, red("Error:"), red("token"), red("framework $PATH"))
		os.Exit(1)
	}

	// Retrieve first argument
	command := strings.ToLower(os.Args[1])

	// Check if command is supported
	switch command {
	case "token":
		GenerateToken()

	case "framework":
		if len(os.Args) < 3 {
			fmt.Printf("Please provide a valid path to generate the framework.\n")
			os.Exit(1)
		}

		path := os.Args[2]
		absPath, err := filepath.Abs(path)
		if err != nil {
			fmt.Printf("\n%s Error processing path: %s\n\n", red("Error:"), red(err.Error()))
			os.Exit(1)
		}

		absPath = filepath.Join(absPath, "")

		GenerateFramework(absPath)
	default:
		// Print error message
		fmt.Printf("\n%s gen CLI tool currently only supports the '%s' command for token generation. \n\n", red("Error:"), red("token"))
		os.Exit(1)
	}
}
