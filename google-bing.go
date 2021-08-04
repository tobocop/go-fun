package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const  (
	BingBaseUrl = "https://www.bing.com/search?q="
	GoogleBaseUrl = "https://www.google.com/search?q="
)

func main() {
	searchTerm := url.QueryEscape(strings.Join(os.Args[1:], ""))

	bingResults := make(chan *http.Response)
	googleResults := make(chan *http.Response)
	errors := make(chan error)

	cx, cancel := context.WithCancel(context.Background())
	bingReq, err := http.NewRequestWithContext(cx, http.MethodGet, fmt.Sprintf("%s%s", BingBaseUrl, searchTerm), nil)
	if err != nil {
		fmt.Printf("Error creating bing request. Err: %v", err)
		os.Exit(1)
	}

	googleRequest, err := http.NewRequestWithContext(cx, http.MethodGet, fmt.Sprintf("%s%s", GoogleBaseUrl, searchTerm), nil)
	if err != nil {
		fmt.Printf("Error creating google request. Err: %v", err)
		os.Exit(1)
	}

	go func() {
		res, err := http.DefaultClient.Do(bingReq)
		if err != nil {
			errors <- err
		}
		bingResults <- res
	}()
	go func() {
		res, err := http.DefaultClient.Do(googleRequest)
		if err != nil {
			errors <- err
		}
		googleResults <- res
	}()

	select {
	case res := <-bingResults:
		cancel()
		defer res.Body.Close()
		err := outputWinner("Bing", res)
		if err != nil {
			fmt.Printf("Error outputting winner. Error %v", err)
			os.Exit(1)
		}
	case res := <-googleResults:
		cancel()
		defer res.Body.Close()
		err := outputWinner("Google", res)
		if err != nil {
			fmt.Printf("Error outputting winner. Error %v", err)
			os.Exit(1)
		}
	case <-errors:
		cancel()
		fmt.Println("One call errored, results irrelevant")
		os.Exit(1)
	}
}

func outputWinner(winner string, res *http.Response) error {
	body := make([]byte, 100)
	_, err := res.Body.Read(body)
	if err != nil {
		return err
	}
	fmt.Println(string(body))
	fmt.Printf("\n%s won\n", winner)
	return nil
}
