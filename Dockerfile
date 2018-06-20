FROM golang:1.10-alpine as builder

COPY . /go/src/github.com/hikhvar/exifsorter
WORKDIR  /go/src/github.com/hikhvar/exifsorter
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /bin/exifsorter .

FROM scratch

COPY --from=builder /bin/exifsorter /bin/exifsorter
ENTRYPOINT ["/bin/exifsorter", "sort", "-s" , "/input", "-t", "/output"]