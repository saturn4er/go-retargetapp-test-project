package main

import (
	"net/http"
	"./topWords"
	"./tcpWordsReceiver"
	"fmt"
	"encoding/json"
	"log"
	"strconv"
	"net"
	"bufio"
	"io"
)

const HTTP_PORT = "8000"
const TCP_PORT = 9000

var provider = topWords.GetTopWordsProvider()

func main() {
	wordsReceiver := tcpWordsReceiver.NewTCPWordsReceiver(TCP_PORT)
	wordsReceiver.OnMessage(provider.AddWordsString)
	wordsReceiver.Start()
	http.HandleFunc("/", wordsResponder)
	http.ListenAndServe(":" + HTTP_PORT, nil)
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