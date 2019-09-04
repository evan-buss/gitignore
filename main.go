package main

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/chzyer/readline"
)

var langFilePath string
var completer *readline.PrefixCompleter

// Find the correct cache path and initialized language completer
func init() {
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal("Could not locate cache directory")
	}
	langFilePath = cacheDir + "/gitignore_languages.txt"

	// Autocomplete with available languages
	completer = readline.NewPrefixCompleter(
		readline.PcItem("refresh"),
		readline.PcItem("create",
			readline.PcItemDynamic(availableLangs(langFilePath),
				readline.PcItemDynamic(availableLangs(langFilePath),
					readline.PcItemDynamic(availableLangs(langFilePath),
						readline.PcItemDynamic(availableLangs(langFilePath),
							readline.PcItemDynamic(availableLangs(langFilePath),
								readline.PcItemDynamic(availableLangs(langFilePath)))))))),

		readline.PcItem("modify",
			readline.PcItem("append",
				readline.PcItemDynamic(availableLangs(langFilePath),
					readline.PcItemDynamic(availableLangs(langFilePath)))),
			readline.PcItem("delete",
				readline.PcItemDynamic(currentLangs(".gitignore"),
					readline.PcItemDynamic(currentLangs(".gitignore"),
						readline.PcItemDynamic(currentLangs(".gitignore"))))),
			readline.PcItem("refresh")),
		readline.PcItem("remove"))
}

func main() {
	checkForLanguageCache(langFilePath)
	gitignoreFile := checkForGitignore(".gitignore")

	// gitignoreFile := &gitignore{langages: currentLangs(:)}

	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[36mgitignore\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		// HistorySearchFold:   true,
		// FuncFilterInputRune: filterInput,
	})
	if err != nil {
		panic(err)
	}
	defer l.Close()

	for {
		line, err := l.Readline()
		if err == readline.ErrInterrupt {
			if len(line) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		parseArguments(line, &gitignoreFile)
	}
}

// checkForLanguageCache checks if the ~/.cache/gitignore_languages.txt exists
// If not, it redownloads it. This file is parsed for language autocompletion
func checkForLanguageCache(path string) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			refreshLanguages(path)
		}
	}
}

// Load internal gitignore object from existing file or create new empty object
func checkForGitignore(path string) Gitignore {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			return Gitignore{}
		}
	}
	body, _ := ioutil.ReadFile(".gitignore")
	return Gitignore{
		languages: currentLangs(".gitignore")(""),
		content:   string(body),
	}
}

// availableLangs parses the languages file to enable language autocompletion
func availableLangs(path string) func(string) []string {
	return func(line string) []string {
		content, err := ioutil.ReadFile(path)
		if err != nil {
			return nil
		}
		return strings.Split(string(content), ",")
	}
}

// currentLangs parses the existing .gitignore languages
// Searches for language comments inside the file "### Java ###"
// Sometimes languages have adittional patches "### Java Patch ###"; ignore these
func currentLangs(path string) func(string) []string {
	return func(line string) []string {
		file, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer file.Close()

		langs := make([]string, 0)
		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			temp := scanner.Text()
			if strings.HasPrefix(temp, "###") &&
				strings.HasSuffix(temp, "###") &&
				!strings.Contains(temp, "Patch") {
				temp = strings.ReplaceAll(temp, "#", "")
				temp = strings.TrimSpace(temp)
				temp = strings.ToLower(temp)
				langs = append(langs, temp)
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		return langs
	}
}

// refreshLanguages downloads a fresh copy of the available languages
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
		log.Fatal("Unable to create new file.")
	}
}
