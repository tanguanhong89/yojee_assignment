FROM golang:1.14

WORKDIR /go/src/github.com/myProj
COPY ./*.go /go/src/github.com/myProj/
COPY ./*.csv /go/src/github.com/myProj/

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build .
RUN ls

ENTRYPOINT ["myProj"]