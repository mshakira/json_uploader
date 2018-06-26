package s3store_test // we are creating a different namespace here

import (
	"fmt"
	"os"
	"testing"

	// our pkg
	"store/s3store"
)

var s3st *s3store.S3Store // this concurrency safe

func TestMain(m *testing.M) {
	var err error
	s3st, err = s3store.Init("jsonuploader-testbucket", "us-west-1")
	//s3st, err = s3store.Init("jsonuploader", "us-west-1")

	if err != nil {
		fmt.Printf("%s Failed, err:%s\n", "s3store.Init", err)
		os.Exit(1)
	}
	var status int = m.Run()

	os.Exit(status)
}

func TestUploadToStore(t *testing.T) {
	var err error
	var payload string = `{"TS" : "2015-09-15T14:00:13Z", "Hello" : "World"}`
	_, err = s3st.UploadToStore("foobar", "key", []byte(payload))
	if err != nil {
		t.Errorf("For %s Expected no error, Error=%s", payload, err)
	}
}
