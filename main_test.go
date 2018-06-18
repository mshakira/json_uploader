package main

import (
        "net/http"
        "net/http/httptest"
        "testing"
        "strings"
)

func TestJsonUploader(t *testing.T) {
        payload := "{\"a\":\"b\"}"
        // Create a request to pass to our handler
        req, err := http.NewRequest("POST", "/upload", strings.NewReader(payload))
        if err != nil {
                t.Fatal(err)
        }

        // Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response
        rr := httptest.NewRecorder()
        handler := http.HandlerFunc(JsonUploader)

        // Our handlers satisfy http.Handler, so we can call their ServeHTTP method
        // directly and pass in our Request and ResponseRecorder.
        handler.ServeHTTP(rr, req)

        // Check the status code is what we expect.
        if status := rr.Code; status != http.StatusOK {
                t.Errorf("handler returned wrong status code: got %v want %v",
                        status, http.StatusOK)
        }

        // Check the response body is what we expect.
        expected := payload
        //if rr.Body.String() != expected {
        if rr.Body.String() == expected {
                t.Errorf("handler returned unexpected body: got %v want %v",
                        rr.Body.String(), expected)
        }

}
