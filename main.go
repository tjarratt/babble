package babble

import (
	"flag"
	"fmt"
	"os"
)

var (
	separator     string
	numberOfWords int
)

func init() {
	flag.IntVar(&numberOfWords, "n", 3, "the number of random words to join")
	flag.StringVar(&separator, "s", "-", "a separator to use when joining words")
}

// TODO: break the random word functionality into windows && unix helpers
func main() {
	if len(os.Args) > 1 {
		checkUsage()
	}

	flag.Parse()
	babbler := NewBabbler()
	babbler.Count = numberOfWords
	babbler.Separator = separator

	println(babbler.Babble())
	return
}

func checkUsage() {
	if os.Args[1] == "help" || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf(`
usage: babble -s [separator] -n [number-of-words]
eg: babble -s="-" -n=5 # holy-moly-guacamole-oily-strombole

The separator between words defaults to '-'
The number of words printed defaults to 3
`)
		os.Exit(1)
	}
}
