package wake

import (
	"log"
	"net/http"
	"runtime"
	"time"
)

//Wake takes a slice of heroku app prefixes and creates a slice of urls from them.
//If it is not the correct time, it continues. If it is the correct time, it gets the urls
func Wake(timezone string, prefixes []string) {
	urls := ConvertPrefixes(prefixes)
	for range time.Tick(30 * time.Minute) {
		ok := IsWakeTime(time.Now(), timezone)
		if !ok {
			continue
		}
		GetUrls(urls)
	}
}

//isWakeTime checks if it is between 0600 and 1800 hours.
func IsWakeTime(t time.Time, timezone string) bool {
	tz, err := time.LoadLocation(timezone)
	if err != nil {
		log.Fatal(err)
	}
	h := t.In(tz).Hour()
	return h >= 6 && h <= 21
}

func ConvertPrefixes(prefixes []string) []string {
	var urls []string
	for _, pre := range prefixes {
		urls = append(urls, "https://"+pre+".herokuapp.com")
	}
	return urls
}

func GetUrls(urls []string) ([]int, error) {
	sc := []int{}
	for _, uri := range urls {
		resp, err := http.Get(uri)
		if err != nil {
			return sc, err
		}
		resp.Body.Close()
		sc = append(sc, resp.StatusCode)
		// io.Copy(ioutil.Discard, resp.Body)
	}
	return sc, nil
}

func GetUrlsConc(urls []string) ([]int, error) {
	g := runtime.GOMAXPROCS(0)
	c := make(chan int, g)
	errc := make(chan error, 1)
	go func() {
		defer close(c)
		for _, url := range urls {
			resp, err := http.Get(url)
			if err != nil {
				errc <- err
				break
			}
			c <- resp.StatusCode
		}
	}()
	sc := []int{}
	for code := range c {
		sc = append(sc, code)
	}
	select {
	case err := <-errc:
		return sc, err
	case <-c:
		return sc, nil
	}
}
