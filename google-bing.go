package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	searchTerm := os.Args[1:]
	bingResults := make(chan *http.Response)
	googleResults := make(chan *http.Response)

	go func() {
		bing, _ := http.Get(fmt.Sprintf("https://www.bing.com/search?q=%s", searchTerm))
		bingResults <- bing
	}()
	go func() {
		google, _ := http.Get(fmt.Sprintf("https://www.google.com/search?q=%s", searchTerm))
		googleResults <- google
	}()

	select {
	case res := <-bingResults:
		outputWinner("Bing", res)
	case res := <-googleResults:
		outputWinner("Google", res)
	}
}

func outputWinner(winner string, res *http.Response) {
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body[:200]))
	fmt.Printf("\n%s won\n", winner)
}
