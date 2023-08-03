package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chaos-io/hystrix-go/hystrix"
)

func main() {
	fmt.Println("start the serve...")
	http.ListenAndServe(":8090", &Handle{})
}

type Handle struct{}

func (h *Handle) ServeHTTP(r http.ResponseWriter, request *http.Request) {
	h.Common(r, request)
}

func (h *Handle) Common(r http.ResponseWriter, request *http.Request) {
	hystrix.ConfigureCommand("mycommand", hystrix.CommandConfig{
		Timeout:                int(3 * time.Second),
		MaxConcurrentRequests:  5,
		SleepWindow:            5000,
		RequestVolumeThreshold: 3,
		ErrorPercentThreshold:  3,
	})

	msg := "success"

	_ = hystrix.Do(
		"mycommand",
		func() error {
			_, err := http.Get("https://www.baidu.com")
			if err != nil {
				fmt.Printf("http get error: %v", err)
				return err
			}
			return nil
		},
		func(err error) error {
			fmt.Printf("handle error: %v\n", err)
			msg = "error"
			return nil
		},
	)
	r.Write([]byte(msg))
}
