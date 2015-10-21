FROM golang:1.5.1

MAINTAINER Harrison Shoebridge <harrison@theshoebridges.com>

RUN go get github.com/bigroom/zombies
RUN go get github.com/gorilla/websocket

ADD . /go/src/github.com/bigroom/vision
ADD ./models/ /go/src/github.com/bigroom/vision/models
ADD ./tunnel/ /go/src/github.com/bigroom/vision/tunnel

WORKDIR /go/src/github.com/bigroom/vision

CMD go run main.go
