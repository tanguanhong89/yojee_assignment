FROM golang:1.14

WORKDIR /go/src/github.com/myProj
COPY . /go/src/github.com/myProj

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build .

CMD ['yojee']