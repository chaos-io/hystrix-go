package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/chaos-io/hystrix-go/hystrix"
)

func main() {
	// easyrun1()
	easyrun2()
}

func easyrun1() {
	err := hystrix.Do(
		"test1",
		func() error {
			// talk to other services
			resp, err := http.Get("https://www.google.com/")
			if err != nil {
				fmt.Printf("get error: %v\n", err)
				return err
			}
			fmt.Printf("get response: %+v\n", resp)
			return nil
		},
		func(err error) error {
			fmt.Printf("get fallback error: %v\n", err)
			return nil
		},
	)
	if err != nil {
		fmt.Printf("do error: %v\n", err)
	}
}

func easyrun2() {
	hystrix.ConfigureCommand("test2", hystrix.CommandConfig{
		Timeout:                int(3 * time.Second),
		MaxConcurrentRequests:  10,
		SleepWindow:            5000, // 当熔断器被打开后，SleepWindow 的时间就是控制过多久后去尝试服务是否可用了。
		RequestVolumeThreshold: 10,   // 一个统计窗口 10 秒内请求数量。达到这个请求数量后才去判断是否要开启熔断
		ErrorPercentThreshold:  30,   // 错误百分比，请求数量大于等于 RequestVolumeThreshold 并且错误率到达这个百分比后就会启动熔断
	})

	err := hystrix.Do("test3", func() error {
		resp, err := http.Get("https://www.google.com/")
		if err != nil {
			fmt.Printf("get error: %v\n", err)
			return err
		}
		fmt.Printf("get response: %+v\n", resp)
		return nil
	}, func(err error) error {
		fmt.Printf("get fallback error: %v\n", err)
		return nil
	})
	if err != nil {
		fmt.Printf("do error: %v\n", err)
	}
}

func easyrun3() {

}
