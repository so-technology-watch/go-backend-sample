FROM golang:1.9-alpine

RUN apk add -U make git gcc libc-dev

WORKDIR $GOPATH/src/go-backend-sample/
ADD . $GOPATH/src/go-backend-sample/

RUN make -f Makefile all && apk del make git gcc libc-dev && \
  rm -rf /gopath/pkg && \
  rm -rf /gopath/src && \
  rm -rf /var/cache/apk/*

EXPOSE 8020
ENTRYPOINT ["/go/bin/todolist"]
