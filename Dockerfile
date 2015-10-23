FROM golang:1.5.1

MAINTAINER Harrison Shoebridge <harrison@theshoebridges.com>

RUN go get github.com/bigroom/zombies
RUN go get github.com/gorilla/websocket

RUN go get github.com/codegangsta/gin

# ADD . /go/src/github.com/bigroom/vision
# ADD ./models/ /go/src/github.com/bigroom/vision/models
# ADD ./tunnel/ /go/src/github.com/bigroom/vision/tunnel

WORKDIR /go/src/github.com/bigroom/vision

CMD gin -i -a=8080 -b="vision"
