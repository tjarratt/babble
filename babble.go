package babble

import (
	"math/rand"
	"strings"
	"sync"
	"time"
)

type Babbler struct {
	Count     int
	Separator string
	Words     []string
	mu        sync.Mutex
	rand      *rand.Rand
}

func NewBabbler() (b Babbler) {
	return NewBabblerWithRand(rand.New(rand.NewSource(time.Now().UnixNano())))
}

func NewBabblerWithRand(rnd *rand.Rand) Babbler {
	return Babbler{
		Count:     2,
		Separator: "-",
		Words:     readAvailableDictionary(),
		rand:      rnd,
	}
}

func (this Babbler) Babble() string {
	pieces := []string{}
	this.mu.Lock()
	for i := 0; i < this.Count; i++ {
		pieces = append(pieces, this.Words[this.rand.Int()%len(this.Words)])
	}
	this.mu.Unlock()
	return strings.Join(pieces, this.Separator)
}
