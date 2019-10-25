package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

// parseArguments takes the user's input and responds appropriately
func parseArguments(line string, gitignoreFile *Gitignore) {
	line = strings.TrimSpace(strings.ToLower(line))

	switch {
	case line == "refresh":
		refreshLanguages(langFilePath)
	case strings.HasPrefix(line, "create"):
		line := strings.TrimSpace(line[6:])
		gitignoreFile.languages = strings.Split(line, " ")
		if len(line) > 0 {
			gitignoreFile.createGitignore()
		} else {
			fmt.Println("Please select one or more languages.")
		}
	case strings.HasPrefix(line, "modify"):
		line := strings.TrimSpace(line[6:])
		switch {
		case strings.HasPrefix(line, "append"):
			line = strings.TrimSpace(line[6:])
			err := os.Remove(".gitignore")
			if err != nil {
				log.Println(err)
			}
			fmt.Println("Old: " + strings.ToLower(strings.Join(gitignoreFile.languages, " ")))
			gitignoreFile.languages = append(gitignoreFile.languages, strings.Split(line, " ")...)
			fmt.Println("New: " + strings.ToLower(strings.Join(gitignoreFile.languages, " ")))
			gitignoreFile.createGitignore()
		case strings.HasPrefix(line, "delete"):
			line = strings.TrimSpace(line[6:])
			gitignoreFile.getDifference(line)
			os.Remove(".gitignore")
			gitignoreFile.createGitignore()
		case strings.HasPrefix(line, "refresh"):
			os.Remove(".gitignore") // remove existing file but keep details to regen
			fmt.Println("Refreshing Languages: " + strings.Join(gitignoreFile.languages, " "))
			gitignoreFile.createGitignore()
		}
	case line == "remove":
		gitignoreFile.remove()
	case line == "exit":
		os.Exit(0)
	case line == "help":
		fallthrough
	default:
		usage()
	}
}

// usage displays program instructions
func usage() {
	fmt.Println(`commands:
	refresh -- update cache of available languages
	create -- create new gitignore file (autocomplete with tab)
	modify -- update an existing gitignore file
	├── append -- add new languages to gitignore (autocomplete with tab)
	├── delete -- remove languages from gitignore (autocomplete with tab)
	├── refresh -- refresh contents of gitignore 
	remove -- delete gitignore
	exit -- exit gitignore`)
}
