package main

import (
	"fmt"
	"net/http"
	"time"
)

type Hours []int

var wakeHours = Hours{6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21}
var halfHour time.Duration = 5 * time.Second

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

//Wake takes a url as a string and a channel string. It checks if it is between 0600 and 1800 hours.
//If it is not the correct time, it waits half an hour before checking again. If it is the correct time, it gets the url and then waits half an hour
//If there was an error, it tries again after five seconds
func Wake(url string, c chan string) {
	ok := isWakeTime(time.Now())
	if ok {
		_, err := http.Get(url)
		if err != nil {
			fmt.Println("failed to hit", url)
			time.Sleep(5 * time.Second)
			go Wake(url, c)
		}
		fmt.Println("successfully hit", url)
		c <- url
		time.Sleep(halfHour)
		go Wake(url, c)
	} else {
		fmt.Println("not the correct time, sleeping for half an hour")
		time.Sleep(halfHour)
		go Wake(url, c)
	}
}

func main() {
	n := []string{"life-together-calculator", "truthify"}
	c := make(chan string)
	for _, pre := range n {
		url := "https://" + pre + ".herokuapp.com"
		go Wake(url, c)
		var wokeDynos []string
		wokeDynos = append(wokeDynos, <-c)
	}
}
