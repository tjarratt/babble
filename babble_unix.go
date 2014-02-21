package main

import (
	"io/ioutil"
	"os"
	"strings"
)

func readAvailableDictionary() (words []string, err error) {
	file, err := os.Open("/usr/share/dict/words")
	if err != nil {
		return
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return
	}

	words = strings.Split(string(bytes), "\n")
	return
}
