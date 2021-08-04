package main

import (
	"fmt"
	"net/http"
)

func main() {
	bingResults := make(chan *http.Response)
	googleResults := make(chan *http.Response)

	searchTerm := "some-search"
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
		fmt.Println(res)
		fmt.Println("Bing won")
	case res := <-googleResults:
		fmt.Println(res)
		fmt.Println("Google won")
	}
}


