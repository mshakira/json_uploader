package jsonfmt_test // we are creating a different namespace here

import (
	"fmt"
	"os"
	"testing"

	// our pkg
	"loadfmt/jsonfmt"
)

var jfmt *jsonfmt.JSONDataFmt // this concurrency safe

func TestMain(m *testing.M) {
	var err error
	jfmt, err = jsonfmt.Init(struct{}{})
	jfmt.Test = true // this is for testing
	if err != nil {
		fmt.Printf("%s Failed, err:%s\n", "jsonfmt.Init", err)
		os.Exit(1)
	}
	var status int = m.Run()

	os.Exit(status)
}

func TestUploadPayload(t *testing.T) {
	var err error
	var payload string = `{"TS" : "2015-09-15T14:00:13Z", "Hello" : "World"}`
	err = jfmt.UploadPayload([]byte(payload))
	if err != nil {
		t.Errorf("For %s Expected no error, Error=%s", payload, err)
	}
}
