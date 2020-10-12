package main

import (
	"fmt"
	"time"
)

//set up times for dyno waker to fire (6am to 9pm EST)
type Hours []int
type WakeTime struct {
	AM Hours
	PM Hours
}

var amHours = Hours{6, 7, 8, 9, 10, 11}
var pmHours = Hours{12, 1, 2, 3, 4, 5, 6, 7, 8, 9}
var wt = WakeTime{amHours, pmHours}

//determine whether it's the correct time to fire
func isWakeTime(t time.Time) bool {
	loc, err := time.LoadLocation("America/New_York")
	if err != nil {
		fmt.Println(err)
		return false
	}
	t = t.In(loc)
	fmt.Println(t.Location(), t.Format("15:04"))
	return false
}


//wake function
func wake() {}

//use goroutines and channels and loop through (use twitter-bot for reference)

func main() {
	ok := isWakeTime(time.Now())
	if ok {
		fmt.Print("Correct time to wake dynos!")
	}
}
