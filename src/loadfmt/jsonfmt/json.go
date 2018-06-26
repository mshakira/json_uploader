package jsonfmt

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"reflect"
	"strconv"
	"time"

	"store/s3store"
)

type JSONDataFmt struct {
	Test bool             // this is for testing
	S3st *s3store.S3Store // let the caller fill it up
	// nothing specific
}

func Init(conn interface{}) (*JSONDataFmt, error) {
	return &JSONDataFmt{}, nil
}

func Destroy(jfmt *JSONDataFmt) error {
	return nil
}

// * Checks whether this is a valid JSON
// * Checks has TS Field
// * Parse the TS
// * Upload to the store
func (jfmt *JSONDataFmt) UploadPayload(data []byte) (err error) {
	// initialize empty interface
	// empty interface will hold value of any type. 
	//Since the json struct is unknown, empty interface is used here
	var payload interface{}

	// check whether JSON is Valid
	payload, err = isJSON(data)
	if err != nil {
		return err
	}

	// check for TS
	var ts string
	ts, err = hasTS(payload)
	if err != nil {
		return errors.New(fmt.Sprintf("Payload does NOT contain TS key (err=%s)", err))
	}

	// Parse TS
	var _tm, tm time.Time
	_tm, err = toTime(ts)
	if err != nil {
		return errors.New(fmt.Sprintf("Not able to convert [%s] to time.Time (err=%s)", ts, err))
	}
	//returns the result of rounding t down to a multiple of d (since the zero time)
	tm = _tm.Truncate(time.Duration(5) * time.Minute) // 5 min
	// prefix YYYYMMDDHHMm with MM rotated every 5 mins
	var prefix string = fmt.Sprintf("%s%s%s%s%s", prefix0(tm.Year()), prefix0(int(tm.Month())), prefix0(tm.Day()), prefix0(tm.Hour()), prefix0(tm.Minute()))
	var key string = fmt.Sprintf("%s%s%s%s%s%s", prefix0(_tm.Year()), prefix0(int(_tm.Month())), prefix0(_tm.Day()),
		prefix0(_tm.Hour()), prefix0(_tm.Minute()), prefix0(_tm.Second()))

	// if this is test, return
	// TODO: write to test bucket?
	if jfmt.Test {
		return nil
	}

	log.Printf("Uploading to Prefix=%s Key=%s\n", prefix, key)
	_, err = jfmt.S3st.UploadToStore(prefix, key, data)
	if err != nil {
		return err
	}

	return err
}

// checks whether the given payload is JSON
// Returns: False if not a JSON, True if valid JSON
func isJSON(payload []byte) (interface{}, error) {
	var p interface{}
	var err error
	//decode json
	err = json.Unmarshal(payload, &p)
	return p, err
}

// check whether we have TS in the payload
// Returns: False if there is not TS field, else True
func hasTS(p interface{}) (string, error) {
	// it should be a Map
	var pVal reflect.Value = reflect.ValueOf(p) //returns type of runtime data
	if pVal.Kind() == reflect.Map {
		// check for TS
		pTSVal := pVal.MapIndex(reflect.ValueOf("TS"))
		if !pTSVal.IsValid() {
			return "", errors.New(fmt.Sprintf("TS field is not present"))
		} else {
			return reflect.ValueOf(pTSVal.Interface()).String(), nil
		}
	} else {
		return "", errors.New(fmt.Sprintf("Payload is not a Map"))
	}
}

// check whether the TS is valid
// Valid TS is defined as below
// * if it can be parsed as RFC3339
func toTime(ts string) (t time.Time, err error) {
	return time.Parse(time.RFC3339, ts)
}

// prefix 0 if 0 is missing and len of string < 2
func prefix0(needle int) (zstr string) {
	var str string = strconv.Itoa(needle)
	if len(str) < 2 {
		zstr = "0" + str
	} else {
		zstr = str
	}
	return zstr
}
