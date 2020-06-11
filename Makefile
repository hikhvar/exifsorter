
.PHONY: test


build: bin/exifsorter

GOVARIABLES=GO111MODULE=on

clean:
	rm -rf bin/exifsorter

bin/exifsorter:
	$(GOVARIABLES) go build -o bin/exifsorter .

test:
	$(GOVARIABLES) go test -cover -race ./...

vet:
	$(GOVARIABLES) go vet ./...

container:
	docker build -t exifsorter .

vendor:
	$(GOVARIABLES) go mod tidy
	$(GOVARIABLES) go mod vendor
