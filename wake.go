package dyno_waker

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

//Wake takes a slice of heroku app prefixes and creates a slice of urls from them.
//If it is not the correct time, it continues. If it is the correct time, it gets the urls
func Wake(timezone string, prefixes []string) {
	urls := convertPrefixes(prefixes)
	ok := IsWakeTime(timezone)
	for range time.Tick(halfHour) {
		if !ok {
			continue
		}
		getUrls(urls)
	}
}

// var wakeHours = Hours{6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21}
var halfHour time.Duration = 30 * time.Minute

//isWakeTime checks if it is between 0600 and 1800 hours.
func IsWakeTime(timezone string) bool {
	ny, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatal(err)
	}
	h := time.Now().In(ny).Hour()
	return h > 6 && h < 21
}

func convertPrefixes(prefixes []string) []string {
	var urls []string
	for _, pre := range prefixes {
		urls = append(urls, "https://"+pre+".herokuapp.com")
	}
	return urls
}

func getUrls(urls []string) {
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
