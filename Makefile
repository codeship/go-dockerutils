.PHONY: \
	all \
	precommit \
	check_for_jet \
	deps \
	updatedeps \
	testdeps \
	updatetestdeps \
	generate \
	build \
	install \
	cov \
	test \
	jetsteps \
	doc \
	clean

all: test install

precommit: doc

check_for_jet:
	@ if ! which jet > /dev/null; then \
		echo "error: jet not installed" >&2; \
		exit 1; \
	  fi

deps:
	go get -d -v ./...

updatedeps:
	go get -d -v -u -f ./...

testdeps: deps
	go get -d -v -t ./...

updatetestdeps: updatedeps
	go get -d -v -t -u -f ./...

generate:
	go generate ./...

build: deps generate
	go build ./...

install: deps generate
	go install ./...

cov: testdeps generate
	go get -v github.com/axw/gocov/gocov
	gocov test | gocov report

test: testdeps generate
	go test -test.v ./...

jetsteps: check_for_jet generate
	jet steps

doc: generate
	go get -v github.com/robertkrimen/godocdown/godocdown
	cp .readme.header README.md
	godocdown | tail -n +7 >> README.md

clean:
	go clean -i ./...
