
.PHONY: test

all: clean container

build: exifsorter

clean:
	rm -rf exifsorter

exifsorter:
	$(GOVARIABLES) go build -o exifsorter .

test:
	$(GOVARIABLES) go test -cover -race ./...

vet:
	$(GOVARIABLES) go vet ./...

lint:
	golangci-lint run

container: build
	docker build -t exifsorter .

vendor:
	$(GOVARIABLES) go mod tidy
	$(GOVARIABLES) go mod vendor
