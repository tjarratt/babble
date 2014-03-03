package babble

import (
	"math/rand"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Babbler struct {
	Count int
	Separator string
	Words []string
}

func NewBabbler() (b Babbler) {
	b.Count = 2
	b.Separator = "-"
	b.Words = readAvailableDictionary()
	return
}

func (this Babbler) Babble() string {
	pieces := []string{}
	for i := 0; i < this.Count ; i++ {
		pieces = append(pieces, this.Words[rand.Int()%len(this.Words)])
	}

	return strings.Join(pieces, this.Separator)
}
