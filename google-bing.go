package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	searchTerm := url.QueryEscape(strings.Join(os.Args[1:], ""))

	bingResults := make(chan *http.Response)
	googleResults := make(chan *http.Response)
	errors := make(chan error)

	cx, cancel := context.WithCancel(context.Background())
	bingReq, err := http.NewRequestWithContext(cx, http.MethodGet, fmt.Sprintf("https://www.bing.com/search?q=%s", searchTerm), nil)
	if err != nil {
		fmt.Printf("Error creating bing request. Err: %v", err)
		os.Exit(1)
	}

	googleRequest, err := http.NewRequestWithContext(cx, http.MethodGet, fmt.Sprintf("https://www.google.com/search?q=%s", searchTerm), nil)
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
		}
	case res := <-googleResults:
		cancel()
		defer res.Body.Close()
		err := outputWinner("Google", res)
		if err != nil {
			fmt.Printf("Error outputting winner. Error %v", err)
		}
	case <-errors:
		cancel()
		fmt.Println("One call errored, results irrelevant")
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
