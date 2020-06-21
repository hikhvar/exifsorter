FROM scratch

COPY exifsorter /bin/exifsorter
ENTRYPOINT ["/bin/exifsorter", "sort", "-s" , "/input", "-t", "/output"]
