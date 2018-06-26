# jsonuploader

Jsonuploader is a web service that uploads the given JSON payload to s3 if it is a valid JSON.

## Directory Structures
```
src
src/uploader.go
src/store
src/store/genericstore.go
src/store/s3store
src/store/s3store/s3_integration_test.go
src/store/s3store/s3.go
src/loadfmt
src/loadfmt/genericfmt.go
src/loadfmt/jsonfmt
src/loadfmt/jsonfmt/json_test.go
src/loadfmt/jsonfmt/json_integration_test.go
src/loadfmt/jsonfmt/json.go
```
- The code structure consists of 3 components.
  - main package - `src/uploader.go`
  - data store - `src/store`
  - package to load formats - `src/loadfmt` 
## Design
![Code Design](img/design.png)
- The code is designed in such a way that it can be extended if requirement rises.
- At present it supports JSON format objects and it will be uploaded to S3 data store.
- In future, we can add support to more formats like YAML, XML etc by following the src/loadfmt/genericfmt.go layouts.
- Simillarly, for data stores, we can support more stores like MYSQL, or custom store by following the src/store/genericstore.go layout.
## Test
- Unit tests and integration tests are added to both loadfmt and S3 stores. To execute the tests, follow the below steps.
```
$ git clone git@github.com:mshakira/jsonuploader.git
$ cd jsonuploader
$ export GOPATH=$(pwd):$GOPATH
$ cd src/store/s3store/
$ go test
PASS
ok  	store/s3store	0.193s
``` 
