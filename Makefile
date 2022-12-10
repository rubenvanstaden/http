SHELL := /bin/bash

UNIT_TEST_PATH=./...

run:
	go mod tidy -compat=1.17
	gofmt -l -s -w .
	go run .

test.unit:
	go test -count=1 -run=Unit $(UNIT_TEST_PATH) -v
