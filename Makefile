default: test build

test:
	cd test; go test -v
	cd crawler; go test -v

build:
	go get -v
	go build -v -o go-crawler



.PHONY: test build
