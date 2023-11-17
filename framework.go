package main

import (
	"archive/zip"
	"bufio"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"github.com/gosimple/slug"
	"golang.org/x/crypto/nacl/secretbox"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var appName *string
var port *int

func GenerateFramework(path string) {
	file, err := content.Open(fileName)
	if err != nil {
		fmt.Printf("\n%s %s \n\n", red("Error:"), red(err))
		return
	}

	defer func(file fs.File) {
		_ = file.Close()
	}(file)

	// Create a new file to write the embedded zip content
	outputFile, err := os.Create(fileName)
	if err != nil {
		fmt.Printf("\n%s %s \n\n", red("Error:"), red(err))
		return
	}
	defer func(outputFile *os.File) {
		_ = outputFile.Close()
	}(outputFile)

	_, err = io.Copy(outputFile, file)
	if err != nil {
		fmt.Printf("\n%s %s \n\n", red("Error:"), red(err))
		return
	}

	// Defer removal of zipped file
	defer func() {
		_ = os.RemoveAll(fileName)
	}()

	// Delete existing path
	_ = os.RemoveAll(path)

	SetAppName()
	SetPort()

	fmt.Printf("\n\n%s\n\n", blue("Generating rust framework (actix-web + async-graphql + sqlx w/ postgres)..."))

	// Create the directory structure
	err = os.MkdirAll(path, os.ModePerm)
	if err != nil {
		fmt.Printf("\n%s %s \n\n", red("Error:"), red(err))
		return
	}

	// Open the zip file
	zipFile, err := zip.OpenReader("framework.zip")
	if err != nil {
		fmt.Printf("%s unable to extract file.\n\n", red("Error:"))
		return
	}

	// Defer zipped file
	defer func(zipFile *zip.ReadCloser) {
		_ = zipFile.Close()
	}(zipFile)

	// Extract the contents of the zip file
	for _, file := range zipFile.File {
		err = ExtractFile(file, path)
		if err != nil {
			fmt.Printf("Error extracting file %s: %s\n", file.Name, err)
			return
		}
	}

	fmt.Printf("1. Generating files in directory... %s", yellow("Done ✓"))

	// Generate a new random 24-byte nonce.
	var nonce [24]byte
	if _, err = io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	// Generate masterKey
	var masterKey [32]byte
	if _, err = io.ReadFull(rand.Reader, masterKey[:]); err != nil {
		fmt.Printf("\n%s Unable to generate random bytes - %s.\n\n", red("Error:"), err)
		return
	}

	// Master key
	masterKeyBytes := masterKey[:]
	b64MasterKey := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(masterKeyBytes)

	// Tokenized master key
	generatedMasterKey := secretbox.Seal(nonce[:], masterKeyBytes, &nonce, &masterKey)
	tokenizedMasterKey := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(generatedMasterKey)

	// Set slug
	s := slug.Make(*appName)

	// Replace .env file
	replacements := map[string]string{
		"APP_NAME=":                    "APP_NAME=" + s,
		"MASTER_KEY=":                  "MASTER_KEY=" + b64MasterKey,
		"# CONTROLLER_BEARER_TOKEN=\n": "CONTROLLER_BEARER_TOKEN=" + tokenizedMasterKey,
	}

	err = ReplaceInFile(path+"/.env", replacements)
	if err != nil {
		fmt.Printf("\n%s %s.\n\n", red("Error:"), err)
		return
	}

	fmt.Printf("\n2. Generating .env variable... %s", yellow("Done ✓"))

	// Generate a new random 24-byte nonce.
	if _, err = io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	// Generate accessTokenKeySigning
	var accessTokenKeySigning [32]byte
	if _, err = io.ReadFull(rand.Reader, accessTokenKeySigning[:]); err != nil {
		fmt.Printf("\n%s Unable to generate random bytes - %s.\n\n", red("Error:"), err)
		return
	}

	// Generate a new random 24-byte nonce.
	if _, err = io.ReadFull(rand.Reader, nonce[:]); err != nil {
		panic(err)
	}

	// Generate refreshTokenKeySigning
	var refreshTokenKeySigning [32]byte
	if _, err = io.ReadFull(rand.Reader, refreshTokenKeySigning[:]); err != nil {
		fmt.Printf("\n%s Unable to generate random bytes - %s.\n\n", red("Error:"), err)
		return
	}

	// Access keys
	accessTokenKeySigningBytes := accessTokenKeySigning[:]
	b64AccessTokenKeySigning := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(accessTokenKeySigningBytes)

	// Refresh keys
	refreshTokenKeySigningBytes := refreshTokenKeySigning[:]
	b64RefreshTokenKeySigning := base64.URLEncoding.WithPadding(base64.NoPadding).EncodeToString(refreshTokenKeySigningBytes)

	// Replace config file
	replacements = map[string]string{
		"pub const SERVER_PORT: u16 = 9020; // Set your port here":                       "pub const SERVER_PORT: u16 = " + strconv.Itoa(*port) + ";",
		"slugify!(\"Rust Server\")":                                                      "slugify!(\"" + *appName + "\")",
		"pub const PASETO_ACCESS_TOKEN_KEY_SIGNING: &str = \"\"; // Use generator here":  "pub const PASETO_ACCESS_TOKEN_KEY_SIGNING: &str = \"" + b64AccessTokenKeySigning + "\";",
		"pub const PASETO_REFRESH_TOKEN_KEY_SIGNING: &str = \"\"; // Use generator here": "pub const PASETO_REFRESH_TOKEN_KEY_SIGNING: &str = \"" + b64RefreshTokenKeySigning + "\";",
		"My Server": *appName,
		"my-server": s,
	}

	err = ReplaceInFile(path+"/backend/config/src/lib.rs", replacements)
	if err != nil {
		fmt.Printf("\n%s %s.\n\n", red("Error:"), err)
		return
	}

	fmt.Printf("\n3. Configuring framework (%s)... %s", red("/backend/config/src/lib.rs"), yellow("Done ✓"))

	// Replace templates
	replacements = map[string]string{
		"My Server": *appName,
		"2023":      strconv.Itoa(time.Now().Year()),
	}

	err = ReplaceInFile(path+"/assets/templates/emails/setup/config.html.hbs", replacements)
	if err != nil {
		fmt.Printf("\n%s %s.\n\n", red("Error:"), err)
		return
	}

	fmt.Printf("\n4. Updating html email templates... %s", yellow("Done ✓"))

	// Configure package.json
	replacements = map[string]string{
		"my-server":               s,
		"$my_project_folder_path": path,
	}

	err = ReplaceInFile(path+"/package.json", replacements)
	if err != nil {
		fmt.Printf("\n%s %s.\n\n", red("Error:"), err)
		return
	}

	fmt.Printf("\n5. Configure package.json file... %s", yellow("Done ✓"))

	fmt.Printf("\n\n%s\n\n", blue("Project Generation Complete!"))
	fmt.Printf("%s %s %s\n", yellow("Navigate to"), cyan("'"+path+"/'"), yellow("and update your .env path for your database connection."))
	fmt.Printf("%s %s %s\n", yellow("You can also start browsing through the"), cyan("README.md"), yellow("file before you start coding."))

	fmt.Printf("\n%s\n\n", blue("Happy Coding!"))
}

func SetAppName() {
	fmt.Printf("\n%s ", yellow("Enter your App Name:"))

	reader := bufio.NewReader(os.Stdin)
	name, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\n%s %s \n", red("Error:"), red(err))
		SetAppName()
	}

	name = strings.TrimSpace(name)

	if len(name) < 1 {
		fmt.Printf("\n%s %s \n", red("Error:"), red("Please provide a valid app name."))
		SetAppName()
	}

	appName = &name
}

func SetPort() {
	fmt.Printf("%s ", yellow("Enter Port:"))

	reader := bufio.NewReader(os.Stdin)
	value, err := reader.ReadString('\n')
	if err != nil {
		fmt.Printf("\n%s %s \n\n", red("Error:"), red(err))
		SetAppName()
	}

	p, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || len(strings.TrimSpace(value)) < 4 {
		fmt.Printf("\n%s \n\n", red("Error: Invalid port number."))
		SetPort()
	}

	port = &p
}

func ExtractFile(file *zip.File, targetPath string) error {
	zipFile, err := file.Open()
	if err != nil {
		return err
	}

	defer func(zipFile io.ReadCloser) {
		_ = zipFile.Close()
	}(zipFile)

	extractPath := filepath.Join(targetPath, file.Name)

	if file.FileInfo().IsDir() {
		err = os.MkdirAll(extractPath, os.ModePerm)
		if err != nil {
			return err
		}

		return nil
	}

	err = os.MkdirAll(filepath.Dir(extractPath), os.ModePerm)
	if err != nil {
		return err
	}

	extractFile, err := os.Create(extractPath)
	if err != nil {
		return err
	}

	defer func(extractFile *os.File) {
		_ = extractFile.Close()
	}(extractFile)

	_, err = io.Copy(extractFile, zipFile)
	if err != nil {
		return err
	}

	return nil
}

func ReplaceInFile(filePath string, replacements map[string]string) error {
	// Read the content of the file
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	// Perform the replacements
	newContent := string(content)
	for oldStr, newStr := range replacements {
		newContent = strings.Replace(newContent, oldStr, newStr, -1)
	}

	// Write the modified content back to the file
	err = os.WriteFile(filePath, []byte(newContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
