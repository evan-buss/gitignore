package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
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
						readline.PcItemDynamic(currentLangs(".gitignore")))))),
		readline.PcItem("remove"))
}

func main() {

	ensureLanguageFileExists(langFilePath)

	l, err := readline.NewEx(&readline.Config{
		// Prompt:          "\033[31m»\033[0m ",
		Prompt:          "\033[31mgitignore\033[0m ",
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

		line = strings.TrimSpace(line)
		switch {
		case line == "refresh":
			fmt.Println("Refreshing gitignore languages cache")
			refreshLanguages(langFilePath)
		case strings.HasPrefix(line, "create"):
			line := strings.TrimSpace(line[6:])

			if len(line) > 0 {
				path, err := createGitignore(line)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("File created at " + path)
				}
			}
		case strings.HasPrefix(line, "modify"):
			line := strings.TrimSpace(line[6:])
			switch {
			case strings.HasPrefix(line, "append"):
				line = strings.TrimSpace(line[6:])
				allLangs := strings.Join(currentLangs(".gitignore")(""), " ")
				fmt.Println("Old: " + allLangs)
				fmt.Println("New: " + allLangs + " " + line)
				os.Remove(".gitignore")
				path, err := createGitignore(allLangs + " " + line)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Gitignore updated at " + path)
				}

			case strings.HasPrefix(line, "delete"):
				line = strings.TrimSpace(line[6:])
				diff := getDifference(line)
				err := os.Remove(".gitignore")
				path, err := createGitignore(diff)
				if err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Gitignore updated at " + path)
				}
			}
		case line == "remove":
			err := os.Remove(".gitignore")
			if err != nil {
				fmt.Println("Couldn't find an existing .gitignore file")
			} else {
				fmt.Println("Existing gitignore file removed")
			}
		case line == "help":
			fallthrough
		default:
			// TODO: Upgrade usage clause to be more detailed
			usage(l.Stderr())
		}
	}
}

func usage(w io.Writer) {
	io.WriteString(w, "commands:\n")
	io.WriteString(w, completer.Tree("    "))
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

// getDifference takes a string containing the languages that should be removed
//   and returns which languages should be kept
//  Ex) Gitignore contains "Java Intellij Python" -> removeLangs = "Intellij Java"
//      getDifference returns "Python"
func getDifference(input string) string {
	allLangs := currentLangs(".gitignore")("")
	removeLangs := strings.Split(input, " ")

	sort.Strings(allLangs)
	sort.Strings(removeLangs)

	output := strings.Join(allLangs, " ")

	for _, lang := range removeLangs {
		output = strings.ReplaceAll(output, lang, "")
	}

	output = strings.TrimSpace(output)
	output = strings.ReplaceAll(output, "  ", " ")
	fmt.Println("Old: " + strings.Join(allLangs, " "))
	fmt.Println("New: " + output)

	return output
}

// currentLangs parses the existing .gitignore languages
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
			if strings.HasPrefix(temp, "###") && strings.HasSuffix(temp, "###") && !strings.Contains(temp, "Patch") {
				langs = append(langs, strings.TrimSpace(strings.ReplaceAll(temp, "#", "")))
			}
		}

		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}

		return langs
	}
}

// ensureLanguageFileExists checks if the ~/.cache/gitignore_languages.txt exists
// If not, it redownloads it. This file is parsed for language autocompletion
func ensureLanguageFileExists(path string) {
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			fmt.Println("Downloading avaiable .gitignore " +
				"languages to ~/.cache/gitignore_languages.txt")
			refreshLanguages(path)
		}
	}
}

// refreshLanguages downloads a fresh copy of the available languages
// The language file is saved to ~/.cache/gitignore_languages.txt
func refreshLanguages(path string) {

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

// createGitignore downloads a .gitignore file to the current directory
// languages - space separated list of languages to be included in .gitignore
// Fails if a .gitignore file already exists.
func createGitignore(languages string) (string, error) {
	if _, err := os.Stat(".gitignore"); err != nil {
		if os.IsNotExist(err) {
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

			path, _ := filepath.Abs(".gitignore")

			return path, nil
		}
	}
	return "", errors.New("gitignore file already exists")
}
