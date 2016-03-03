package main

import (
	"net/http"
	"./topWords"
	"fmt"
	"encoding/json"
	"log"
	"strconv"
	"net"
	"os"
	"bufio"
	"io"
)

const HTTP_PORT = "8000"
const TCP_PORT = "9000"

var provider = topWords.GetTopWordsProvider()

func main() {
	go runTCPWordsReceiver()
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
func runTCPWordsReceiver() {
	l, err := net.Listen("tcp", ":" + TCP_PORT)
	if err != nil {
		log.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	log.Println("Waiting words on :" + TCP_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleTCPConnection(conn)
	}
}
func handleTCPConnection(conn net.Conn) {
	defer conn.Close()
	conn_reader := bufio.NewReader(conn)
	for {
		line, _, err := conn_reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				log.Printf("Client %v disconnected", conn.RemoteAddr())
			}else {
				fmt.Println(err)
			}
			break
		}
		provider.AddWordsString(string(line))
	}
}