package dyno_waker

import (
	"reflect"
	"testing"
)

func TestConvertPrefixes(t *testing.T) {
	dynos := []string{"life-together-calculator", "truthify"}
	want := []string{"https://life-together-calculator.herokuapp.com", "https://truthify.herokuapp.com"}
	got := convertPrefixes(dynos)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("want %v; got %v", want, got)
	}
}
