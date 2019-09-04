package main

import (
	"reflect"
	"testing"
)

func TestGetDifference(t *testing.T) {
	file := Gitignore{languages: []string{"java", "intellij", "go"}}
	file.getDifference("go")
	if !reflect.DeepEqual(file.languages, []string{"intellij", "java"}) {
		t.Errorf("getDifference = %s; want [intellij, java]", file.languages)
	}
}
func TestGetDifferenceBug(t *testing.T) {
	file := Gitignore{languages: []string{"java", "intellij", "go"}}
	file.getDifference("intellij java")
	if !reflect.DeepEqual(file.languages, []string{"go"}) {
		t.Errorf("getDifference = %s; want [go]", file.languages)
	}
}
func TestGetDifferenceEmpty(t *testing.T) {
	file := Gitignore{languages: []string{}}
	file.getDifference("go")
	if !reflect.DeepEqual(file.languages, []string{}) {
		t.Errorf("getDifference = %s; want []]", file.languages)
	}
}
