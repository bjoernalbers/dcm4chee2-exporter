package main

import (
	"reflect"
	"testing"
)

func TestToMap(t *testing.T) {
	in := []byte(`
# A comment
MessageCount=3
DeliveringCount=0
ScheduledMessageCount=3

foo=bar
`)
	got := Translate(in)
	want := map[string]int{"MessageCount": 3,
		"DeliveringCount":       0,
		"ScheduledMessageCount": 3,
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("Translate() returned unespected result:\ngot:\t%v\nwant:\t%v", got, want)
	}
}
