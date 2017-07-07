PHONY: test deps check

test:
	go test -v

deps:
	go get -u honnef.co/go/tools/cmd/staticcheck
	go get -u honnef.co/go/tools/cmd/gosimple
	go get -u honnef.co/go/tools/cmd/unused

check:
	staticcheck $$(go list ./... | grep -v /vendor/)
	gosimple $$(go list ./... | grep -v /vendor/)
	unused $$(go list ./... | grep -v /vendor/)

bench:
	go test -v -bench=. 

benchmem:
	go test -v -benchmem -bench=. 