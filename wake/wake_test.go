package wake_test

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/cpustejovsky/dyno-waker/wake"
	"github.com/stretchr/testify/assert"
)

func TestConvertPrefixes(t *testing.T) {
	dynos := []string{"life-together-calculator", "truthify"}
	want := []string{"https://life-together-calculator.herokuapp.com", "https://truthify.herokuapp.com"}
	got := wake.ConvertPrefixes(dynos)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v; got %v", want, got)
	}
}

func TestIsWakeTime(t *testing.T) {
	tests := []struct {
		name     string
		timezone string
		hour     int
		want     bool
	}{
		{"000 hours (EST)", "America/New_York", 0, false},
		{"100 hours (EST)", "America/New_York", 1, false},
		{"200 hours (EST)", "America/New_York", 2, false},
		{"300 hours (EST)", "America/New_York", 3, false},
		{"400 hours (EST)", "America/New_York", 4, false},
		{"500 hours (EST)", "America/New_York", 5, false},
		{"600 hours (EST)", "America/New_York", 6, true},
		{"700 hours (EST)", "America/New_York", 7, true},
		{"800 hours (EST)", "America/New_York", 8, true},
		{"900 hours (EST)", "America/New_York", 9, true},
		{"1000 hours (EST)", "America/New_York", 10, true},
		{"1100 hours (EST)", "America/New_York", 11, true},
		{"1200 hours (EST)", "America/New_York", 12, true},
		{"1300 hours (EST)", "America/New_York", 13, true},
		{"1400 hours (EST)", "America/New_York", 14, true},
		{"1500 hours (EST)", "America/New_York", 15, true},
		{"1600 hours (EST)", "America/New_York", 16, true},
		{"1700 hours (EST)", "America/New_York", 17, true},
		{"1800 hours (EST)", "America/New_York", 18, true},
		{"1900 hours (EST)", "America/New_York", 19, true},
		{"2000 hours (EST)", "America/New_York", 20, true},
		{"2100 hours (EST)", "America/New_York", 21, true},
		{"2200 hours (EST)", "America/New_York", 22, false},
		{"2300 hours (EST)", "America/New_York", 23, false},
		{"2400 hours (EST)", "America/New_York", 24, false},
	}
	for _, tt := range tests {
		tz, err := time.LoadLocation(tt.timezone)
		if err != nil {
			log.Fatal(err)
		}
		testTime := time.Date(2020, 10, 14, tt.hour, 0, 0, 0, tz)
		got := wake.IsWakeTime(testTime, tt.timezone)
		if got != tt.want {
			t.Errorf("want %v; got %v", tt.want, got)
		}
	}

}

func TestGetUrls(t *testing.T) {
	expected := "dummy data"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	}))
	defer svr.Close()
	urls := []string{svr.URL, svr.URL, svr.URL}

	want := []int{http.StatusOK, http.StatusOK, http.StatusOK}

	got, err := wake.GetUrls(urls)
	if err != nil {
		t.Errorf("Got error:\t%v", err)
	}
	assert.Equal(t, want, got)
}

func TestGetUrlsConc(t *testing.T) {
	expected := "dummy data"
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, expected)
	}))
	defer svr.Close()
	urls := []string{svr.URL, svr.URL, svr.URL}

	want := []int{http.StatusOK, http.StatusOK, http.StatusOK}

	got, err := wake.GetUrlsConc(urls)
	if err != nil {
		t.Errorf("Got error:\t%v", err)
	}
	assert.Equal(t, want, got)

}
