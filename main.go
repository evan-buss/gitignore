package main

import (
	"bufio"
	"io"
	"io/ioutil"
	"log"
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
		readline.PcItem("remove"),
		readline.PcItem("exit"))
}

func main() {
	checkForLanguageCache(langFilePath)
	gitignoreFile := checkForGitignore(".gitignore")

	l, err := readline.NewEx(&readline.Config{
		Prompt:          "\033[36mgitignore\033[0m ",
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
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
	body, customRules := parseContent(".gitignore")
	return Gitignore{
		languages:   currentLangs(".gitignore")(""),
		content:     body,
		customRules: customRules,
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

// parseContent reads the gitignore file and separates the
// custom rules from the rest of the body
func parseContent(path string) (string, []string) {
	file, err := os.Open(path)
	if err != nil {
		log.Println("Couldn't parse custom rulesets")
	}
	defer file.Close()

	rules := make([]string, 0)
	var body strings.Builder

	scanner := bufio.NewScanner(file)
	insideCustom := false // We have reached a custom rules block

	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, "# start persist") {
			insideCustom = true
			continue
		} else if strings.Contains(line, "# end persist") {
			insideCustom = false
			continue
		}

		if insideCustom {
			rules = append(rules, line)
		} else {
			body.WriteString(line + "\n")
		}
	}

	return body.String(), rules
}
