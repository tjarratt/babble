// +build darwin dragonfly freebsd linux netbsd openbsd

package babble

import (
	"io/ioutil"
	"os"
	"strings"
)

func readAvailableDictionary() (words []string) {
	file, err := os.Open("/usr/share/dict/words")
	if err != nil {
		panic(err)
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		panic(err)
	}

	words = strings.Split(string(bytes), "\n")
	return
}
