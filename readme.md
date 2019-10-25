# gitignore

A simple tool to manage your project's gitignore files from the command line.

[![asciicast](https://asciinema.org/a/84hL2l0AULc6oxdCUZ1I2hgLE.svg)](https://asciinema.org/a/84hL2l0AULc6oxdCUZ1I2hgLE)

## About

gitignore downloads language specific `.gitignore` files from gitignore.io

## Installation

`go get github.com/evan-buss/gitignore`

## Usage

```
commands:
	refresh -- update cache of available languages
	create -- create new gitignore file (autocomplete with tab)
	modify -- update an existing gitignore file
	├── append -- add new languages to gitignore (autocomplete with tab)
	├── delete -- remove languages from gitignore (autocomplete with tab)
	├── refresh -- refresh contents of gitignore 
	remove -- delete gitignore
```

## Custom Rules

Gitignore downloads language specific rules from gitignore.io. Each time you update or create a new file, the entire
contents are overwritten with the data returned from the API. This is not good if you have specific rules that differ
from the universal ruleset. As a result gitignore has a custom syntax where you can place rules that you don't want to
be overwritten.

Create a new block with the following syntax and your rules will not be overwritten. Put any rules you need in between
the header and footer comment blocks.

```
# start persist
[your rules here]
# done persist
```