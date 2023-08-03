package main

import (
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"
)

var client *http.Client

func init() {
	client = &http.Client{Transport: &http.Transport{
		MaxIdleConns:    100,
		IdleConnTimeout: 1 * time.Second,
	}}
}

func main() {
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			req, err := http.NewRequest("GET", "http://localhost:8090", nil)
			if err != nil {
				fmt.Printf("new request errorr: %v\n", err)
				return
			}

			resp, err := client.Do(req)
			if err != nil {
				fmt.Printf("do request error: %v\n", err)
				return
			}
			defer resp.Body.Close()

			nByte, err := io.ReadAll(resp.Body)
			if err != nil {
				fmt.Printf("get response body error: %v\n", err)
				return
			}
			fmt.Printf("get the message: %v\n", string(nByte))
		}()
	}

	wg.Wait()
	fmt.Printf("DONE\n")
}
