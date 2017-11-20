# minimal linux distribution
FROM golang:1.9-alpine

# GO and PATH env variables already set in golang image
# to reduce download time
RUN apk add -U make git

# set the go path to import the source project
WORKDIR $GOPATH/src/go-backend-sample/
ADD . $GOPATH/src/go-backend-sample/

# In one command-line (for reduce memory usage purposes),
# we install the required software,
# we build bookstore program
# we clean the system from all build dependencies
RUN make -f MakeFile all && apk del make git && \
  rm -rf /gopath/pkg && \
  rm -rf /gopath/src && \
  rm -rf /var/cache/apk/*

# by default, the exposed ports are 8020 (HTTP)
EXPOSE 8020
ENTRYPOINT ["/go/bin/bookstore"]
