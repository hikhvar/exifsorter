
.PHONY: test

build: bin/exifsorter

bin/exifsorter: test vet
	go build -o bin/exifsorter .

test:
	go test -cover -race ./...

vet:
	go vet ./...

container:
	docker build -t exifsorter .