package main

import (
	"net/http"
	"./topWords"
	"fmt"
	"encoding/json"
	"log"
	"strconv"
)

var provider = topWords.GetTopWordsProvider()

func main() {
	http.HandleFunc("/", wordsResponder)
	http.ListenAndServe(":8000", nil)

}
func wordsResponder(w http.ResponseWriter, r *http.Request) {
	topCountS := r.URL.Query()["N"]
	topCount := 10
	if len(topCountS) > 0 {
		_topCount, err := strconv.ParseInt(topCountS[0], 10, 32)
		if err != nil {
			fmt.Println(err)
		} else {
			topCount = int(_topCount)
		}

	}

	topWords := provider.GetTopWords(topCount)
	topWordsBJson, err := json.Marshal(topWords)
	if err != nil {
		log.Panic(err)
	}
	if topWordsBJson == nil {
		topWordsBJson = []byte{'[', ']'}
	}
	fmt.Fprintf(w, `{"top_words":%s}`, topWordsBJson)

}