package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// Gitignore holds details about the current status of the gitignore file
type Gitignore struct {
	languages []string
	content   string
}

// getDifference takes a string containing the languages that should be removed
//   and returns which languages should be kept
//  Ex) Gitignore contains "Java Intellij Python" -> removeLangs = "Intellij Java"
//      getDifference returns "Python"
func (file *Gitignore) getDifference(input string) {
	// allLangs := currentLangs(".gitignore")("")
	allLangs := file.languages
	removeLangs := strings.Split(input, " ")

	sort.Strings(allLangs)
	sort.Strings(removeLangs)

	output := make([]string, 0)
	for _, val := range allLangs {
		if !arrayContains(removeLangs, val) {
			output = append(output, val)
		}
	}
	fmt.Println("Old: " + strings.ToLower(strings.Join(allLangs, " ")))
	fmt.Println("New: " + strings.ToLower(strings.Join(output, " ")))

	file.languages = output
}

// arrayContains determines if the given []string array contains the given value
func arrayContains(array []string, value string) bool {
	for _, val := range array {
		if val == value {
			return true
		}
	}
	return false
}

// createGitignore downloads a .gitignore file to the current directory
// Fails if a .gitignore file already exists.
func (file *Gitignore) createGitignore() {
	if _, err := os.Stat(".gitignore"); err != nil {
		if os.IsNotExist(err) {
			languages := strings.Join(file.languages, " ")
			resp, err := http.Get("https://gitignore.io/api/" +
				strings.ReplaceAll(languages, " ", ","))

			if err != nil {
				log.Fatal("Could not connect to gitignore.io server")
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Could not read response")
			}

			err = ioutil.WriteFile(".gitignore", body, 0644)
			if err != nil {
				log.Fatal("Unable to create new file.")
			}

			file.content = string(body)

			path, _ := filepath.Abs(".gitignore")

			fmt.Println("Gitignore updated at " + path)
			return
		}
	}
	fmt.Println("Error: .gitignore file already exists.")
}

func (file *Gitignore) remove() {

	if err := os.Remove(".gitignore"); err != nil {
		fmt.Println("Couldn't find an existing .gitignore file")
	} else {
		file.languages = nil
		file.content = ""
		fmt.Println("Removed .gitignore file")
	}
}
