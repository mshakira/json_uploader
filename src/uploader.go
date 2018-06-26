package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"

	// our okg
	"loadfmt/jsonfmt"
	"store/s3store"
)

type uploader struct {
	js *jsonfmt.JSONDataFmt
}

var debug bool
var quiet bool
var verbose bool
var help bool
var listener string

// before main
// don't do any "Fatal" stuff here
// go test ./... will fail (it runs 'init')
func init() {
	log.SetFlags(0)
	parseFlags()
	if help == true {
		printHelp()
		os.Exit(0)
	}

	// turn on debuging and force quiet to off
	if verbose {
		debug = true
		quiet = false
	}

}

func parseFlags() {
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.BoolVar(&quiet, "quiet", false, "Enable quiet mode")
	flag.BoolVar(&verbose, "verbose", false, "Enable verbose mode")
	flag.BoolVar(&help, "help", false, "Help")
	flag.StringVar(&listener, "listener", "0.0.0.0:8080", "Interface for HTTP Listener (default=0.0.0.0:8080)")
	flag.Parse()

	return
}

func printHelp() {
	fmt.Printf(
		`Usage: %s[OPTIONS]
eg: %s 
    -debug .................... Debug
    -verbose .................. Verbose Mode (Each Feed Status)
    -quiet .................... Do NOT show even INFO messages
    -help ..................... Good Ol' Help
    -listener ................. HTTP Listener
https://github.com/mshakira/json_uploader
`, os.Args[0], os.Args[0])
	os.Exit(0)
	return
}

// future need to closure the function with more data to be passed?
func genericHandlerV1(fn func(http.ResponseWriter, *http.Request, uploader), u uploader) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) { // returns void, we don't care!
		fn(w, r, u)
	}
}

//handle POST request
func requestHandler(w http.ResponseWriter, r *http.Request, u uploader) {
	if r.Method == "POST" {
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body",
				http.StatusInternalServerError)
		}
		err = u.js.UploadPayload(body)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"status" : "not-uploaded", "error" : "%s"}`, err), http.StatusInternalServerError)
		} else {
			fmt.Fprint(w, `{"status": "uploaded", "error" : ""}`)
		}
	} else {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
	}
}

//to get region value for s3 bucket.
// not used as of now as same bucket is used across the regions. If latency increases, we might revisit
func findRegion() string {
	var (
		cmdOut []byte
		err error
	)
	cmdName := "curl"
	cmdArgs := []string{"-s","http://instance-data/latest/meta-data/placement/availability-zone"}
	if cmdOut, err = exec.Command(cmdName, cmdArgs...).Output(); err != nil {
		fmt.Fprintln(os.Stderr, "Not able to get availability-zone name: ", err)
		os.Exit(1)
	}
	az := string(cmdOut)
	region := az[:len(az)-1]
	return region
}

// for checking if the service is up and running
// usually this path will be monitored by load balancer targer group. If fails, autoscaling group will spawn new instances 
func HeartbeatHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func main() {
	// init the uploader
	var u uploader
	var err error
	// init json
	u.js, err = jsonfmt.Init(struct{}{})
	if err != nil {
		log.Fatal(err)
	}

	//region := findRegion()
	// open the store
	u.js.S3st, err = s3store.Init("jsonuploader", "us-west-1")

	//hosting static web page to upload request
	http.Handle("/", http.FileServer(http.Dir("./static")))

	// handling upload request
	http.HandleFunc("/status.html",HeartbeatHandler)
	http.HandleFunc("/api/v1/upload", genericHandlerV1(requestHandler, u))
	http.HandleFunc("/api/v1/upload/", genericHandlerV1(requestHandler, u)) // handler for trailing

	// start the server
	log.Printf("Server Starting @ [%s]", listener)
	if err := http.ListenAndServe(listener, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
