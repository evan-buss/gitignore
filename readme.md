# gitignore

A simple tool to manage your project's gitignore files from the command line.

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