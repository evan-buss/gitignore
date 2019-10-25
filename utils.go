package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

// arrayContains determines if the given []string array contains the given value
func arrayContains(array []string, value string) bool {
	for _, val := range array {
		if val == value {
			return true
		}
	}
	return false
}

// downloadLanguageList downloads a fresh copy of the available languages
// The language file is saved to ~/.cache/gitignore_languages.txt
func refreshLanguages(path string) {
	fmt.Println("Downloading available languages.")

	resp, err := http.Get("https://gitignore.io/api/list")
	if err != nil {
		log.Fatal("Could not connect to gitignore.io server")
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Could not read response")
	}

	formatted := strings.ReplaceAll(string(body), "\n", ",")

	err = ioutil.WriteFile(path, []byte(formatted), 0644)
	if err != nil {
		log.Fatal(err)
	}
}
