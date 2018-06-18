package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
)

func JsonUploader(w http.ResponseWriter, r *http.Request) {
        if r.Body == nil {
                http.Error(w, "Please send a request body", 400)
                return
        }
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	fmt.Fprintln(w, string(bodyBytes))
}

func main() {
	http.HandleFunc("/upload", JsonUploader)
	if err := http.ListenAndServe(":80", nil); err != nil {
		log.Fatal(err)
	}
}
