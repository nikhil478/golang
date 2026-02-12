### Go Race Detector

-> go provides race detector tool for finding race conditions in go code 

-> go test -race mypkg
-> go run -race mysrc.go
-> go build -race mycmd
-> go install -race mypkg

* Binary needs to be race enabled

* When racy behaviour is detected a warning is printed

* Race enabled binary will 10 times slower and consume 10 times more memory

* Integration tests and load tests are good candidates to test with binary with race enabled