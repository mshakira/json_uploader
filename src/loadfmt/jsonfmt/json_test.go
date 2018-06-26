package jsonfmt

import (
	_ "fmt"
	"testing"
	"time"
)

func TestIsJSON(t *testing.T) {
	// +ve test cases
	var goodTest []string = []string{
		`{"foo":"bar"}`,
		`{"foo":{"bar" : "foobar"} }`,
	}

	// -ve test cases
	var badTest []string = []string{
		`"foo":"bar"`,
		`{"foo":{"bar" : "foobar"}`,
	}

	var err error

	// all is good
	for _, test := range goodTest {
		_, err = isJSON([]byte(test))
		if err != nil {
			t.Errorf("Expected %s to be a valid json, [err=%s]", test, err)
		}
	}

	// all are bad
	for _, test := range badTest {
		_, err = isJSON([]byte(test))
		if err == nil {
			t.Errorf("Expected %s to be a non-valid json, received [err=%#v]", test, err)
		}
	}

	return
}

func TestHasTS(t *testing.T) {
	var err error
	var p interface{}
	var ts string
	var goodTest []string = []string{
		`{"TS":"xyz", "foo":"bar"}`,
		`{"foo":{"bar" : "foobar"}, "TS":"11" }`,
	}
	// all is good
	for _, test := range goodTest {
		p, err = isJSON([]byte(test))

		if err != nil {
			t.Errorf("Expected %s to be a non-valid json, received [isjson=%s]", test, err)
		}
		ts, err = hasTS(p)
		if err != nil || ts == "" {
			t.Errorf("Expected %s to have TS, received [err=%s]", test, err)
		}
	}

	var badTest []string = []string{
		`{"TS1":"xyz", "foo":"bar"}`,
		`{"foo":{"bar" : "foobar"}}`,
	}
	// all is bad
	for _, test := range badTest {
		p, err = isJSON([]byte(test))
		if err != nil {
			t.Errorf("Expected %s to be a non-valid json, received [isjson=%s]", test, err)
		}
		_, err = hasTS(p)
		if err == nil {
			t.Errorf("Expected Error %s NOT to have TS, received [err=%s]", test, err)
		}
	}

}

func TestToTime(t *testing.T) {
	var goodTest []string = []string{
		`2015-09-15T14:00:12-00:00`,
		`2015-09-15T14:00:13Z`,
		`0000-09-15T14:00:13Z`, // this is valid, though year is 0000 (should we put more constraints on valid time)
	}

	var err error
	for _, test := range goodTest {
		_, err = toTime(test)
		if err != nil {
			t.Errorf("Expected %s to be a valid time, received [error=%s]", test, err)
		}
	}

	var badTest []string = []string{
		`0000-01-02T15:04:05ZZ`,
	}
	var tm time.Time
	for _, test := range badTest {
		tm, err = toTime(test)
		if err == nil {
			t.Errorf("Expected %s to be NOT a valid time, received [converted to=%s]", test, tm)
		}
	}
}
