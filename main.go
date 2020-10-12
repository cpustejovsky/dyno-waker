package main

import (
	"fmt"
	"net/http"
	"time"
)

//set up times for dyno waker to fire (6am to 9pm EST)
type Hours []int

var wakeHours = Hours{6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 19, 20, 21}
var halfHour time.Duration = 5 * time.Second

//determine whether it's the correct time to fire
func isWakeTime(t time.Time) bool {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		return false
	}
	t = t.In(loc)
	for _, hr := range wakeHours {
		if t.Hour() == hr {
			return true
		}
	}
	return false
}

//wake function
func wake(url string, c chan string) {
	ok := isWakeTime(time.Now())
	if ok {
		_, err := http.Get(url)
		if err != nil {
			fmt.Println("failed to hit", url)
			time.Sleep(5 * time.Second)
			wake(url, c)
		}
		fmt.Println("successfully hit", url)
		c <- url
	}
	fmt.Println("not the correct time, sleeping for half an hour")
	time.Sleep(halfHour)
	wake(url, c)
}

//use goroutines and channels and loop through (use twitter-bot for reference)

func main() {
	n := []string{"life-together-calculator", "truthify"}
	c := make(chan string)
	for _, pre := range n {
		url := "https://" + pre + ".herokuapp.com"
		go wake(url, c)
		var wokeDynos []string
		wokeDynos = append(wokeDynos, <-c)
	}
}
