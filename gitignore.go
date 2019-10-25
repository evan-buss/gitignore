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

// Gitignore holds the details and contents of the gitignore file
type Gitignore struct {
	languages   []string // Holds the current languages downloaded from API
	content     string   // Holds the entire content body
	customRules []string // Holds the user's custom rulesets that shouldn't be overwritten
}

// getDifference takes a string containing the languages that should be removed
//   and returns which languages should be kept
//  Ex) Gitignore contains "Java Intellij Python" and input contains "IntelliJ Java";
//      getDifference returns "Python"
func (gitignore *Gitignore) getDifference(input string) {
	allLangs := gitignore.languages
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

	gitignore.languages = output
}

// createGitignore downloads a .gitignore file to the current directory
// Fails if a .gitignore file already exists.
func (gitignore *Gitignore) createGitignore() {
	if _, err := os.Stat(".gitignore"); err != nil {
		if os.IsNotExist(err) {
			languages := strings.Join(gitignore.languages, ",")
			resp, err := http.Get("https://gitignore.io/api/" + languages)

			if err != nil {
				log.Fatal("Could not connect to gitignore.io server")
			}

			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Fatal("Could not read response")
			}

			file, err := os.Create(".gitignore")
			if err != nil {
				log.Fatal(err)
			}
			defer file.Close()

			if len(gitignore.customRules) > 0 {
				file.WriteString("# start persist\n")
				for _, line := range gitignore.customRules {
					file.WriteString(line + "\n")
				}
				file.WriteString("# end persist\n")
			}

			_, err = file.Write(body)
			if err != nil {
				log.Fatal(err)
			}

			gitignore.content = string(body)

			path, _ := filepath.Abs(".gitignore")

			fmt.Println("Gitignore updated at " + path)
			return
		}
	}
	fmt.Println("Error: .gitignore file already exists.")
}

func (gitignore *Gitignore) remove() {
	if err := os.Remove(".gitignore"); err != nil {
		fmt.Println("Couldn't find an existing .gitignore file")
	} else {
		gitignore.languages = nil
		gitignore.content = ""
		gitignore.customRules = nil
		fmt.Println("Removed .gitignore file")
	}
}
