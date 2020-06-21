
.PHONY: test

all: clean container

build: exifsorter

GOVARIABLES=GO111MODULE=on

clean:
	rm -rf exifsorter

exifsorter:
	$(GOVARIABLES) go build -o exifsorter .

test:
	$(GOVARIABLES) go test -cover -race ./...

vet:
	$(GOVARIABLES) go vet ./...

container: build
	docker build -t exifsorter .

vendor:
	$(GOVARIABLES) go mod tidy
	$(GOVARIABLES) go mod vendor
