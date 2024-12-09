FROM scratch

COPY exifsorter /bin/exifsorter
ENTRYPOINT ["/bin/exifsorter"]
CMD ["sort", "-s" , "/input", "-t", "/output"]