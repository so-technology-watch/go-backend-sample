FROM golang:1.9-alpine

RUN apk add -U make git

WORKDIR $GOPATH/src/go-backend-sample/
ADD . $GOPATH/src/go-backend-sample/

RUN make -f MakeFile all && apk del make git && \
  rm -rf /gopath/pkg && \
  rm -rf /gopath/src && \
  rm -rf /var/cache/apk/*

EXPOSE 8020
ENTRYPOINT ["/go/bin/todolist"]
