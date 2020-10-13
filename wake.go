package wake

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// type Hours []int

// var wakeHours = Hours{6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21}
var halfHour time.Duration = 30 * time.Minute

//isWakeTime checks if it is between 0600 and 1800 hours.
func isWakeTime() bool {
	ny, err := time.LoadLocation("America/New_York")
	if err != nil {
		log.Fatal(err)
	}
	h := time.Now().In(ny).Hour()
	return h > 6 && h < 21
}

//Wake takes a slice of heroku app prefixes and creates a slice of urls from them.
//If it is not the correct time, it continues. If it is the correct time, it gets the urls
func Wake(prefixes []string) {
	var urls []string
	for _, pre := range prefixes {
		urls = append(urls, "https://"+pre+".herokuapp.com")
	}
	ok := isWakeTime()
	for range time.Tick(halfHour) {
		if !ok {
			continue
		}
		for _, uri := range urls {
			resp, err := http.Get(uri)
			switch {
			case err != nil:
				log.Printf("%s: %v", uri, err)
			case err == nil:
				log.Printf("%s: %d", uri, resp.StatusCode)
				io.Copy(ioutil.Discard, resp.Body)
				resp.Body.Close()
			}
		}
	}
}
