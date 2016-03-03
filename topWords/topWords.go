package topWords

import (
	"sync"
	"sort"
)

type Pair struct {
	Key   string
	Value int
}
type PairList []Pair

func (p PairList) Len() int {
	return len(p)
}
func (p PairList) Less(i, j int) bool {
	return p[i].Value > p[j].Value
}
func (p PairList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
type TopWordsProviderPointer *topWordsProvider
type topWordsProvider struct {
	top              []string
	topMutex         sync.Mutex
	topCount         int
	wordsCounts      map[string]int
	wordsCountsMutex sync.Mutex
}

// Add string to top
func (this *topWordsProvider) AddWordsString(wordsString string) {
	this.wordsCountsMutex.Lock()
	defer this.wordsCountsMutex.Unlock()

	var prevIndex int
	var prevWasSpace bool
	// Search words divided by space ( It faster than regexp )
	for i, c := range wordsString {
		if c == ' ' {
			if prevIndex != -1 {
				this.wordsCounts[wordsString[prevIndex:i]] += 1
				prevIndex = -1
			}
			prevWasSpace = true
			continue
		}else {
			if prevWasSpace {
				prevIndex = i
			}
			prevWasSpace = false
		}

	}
	this.recalculateTop()
}

func (this *topWordsProvider) recalculateTop() {
	this.topMutex.Lock()
	defer this.topMutex.Unlock()
	this.topCount = len(this.wordsCounts)
	pl := make(PairList, this.topCount)
	i := 0
	for k, v := range this.wordsCounts {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(pl)
	newTop := make([]string, this.topCount)
	for i, value := range pl {
		newTop[i] = value.Key
	}
	this.top = newTop
}

func (this *topWordsProvider) GetTopWords(count int) []string {
	this.topMutex.Lock()
	defer this.topMutex.Unlock()
	if count > this.topCount {
		count = this.topCount
	}
	result := make([]string, count)
	copy(result, this.top[:count])
	return result
}

func GetTopWordsProvider() *topWordsProvider {
	result := new(topWordsProvider)
	result.wordsCounts = map[string]int{}
	return result
}