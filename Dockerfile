FROM golang:1.5.1

MAINTAINER Harrison Shoebridge <harrison@theshoebridges.com>

RUN go get github.com/gorilla/websocket
RUN go get github.com/gorilla/mux
RUN go get github.com/paked/gerrycode/communicator
RUN go get github.com/paked/configure
RUN go get github.com/paked/restrict
RUN go get gopkg.in/mgutz/dat.v1
RUN go get gopkg.in/mgutz/dat.v1/sqlx-runner
RUN go get github.com/dgrijalva/jwt-go
RUN go get github.com/sorcix/irc
RUN go get github.com/nickvanw/ircx

RUN go get github.com/codegangsta/gin

# ADD . /go/src/github.com/bigroom/vision
# ADD ./models/ /go/src/github.com/bigroom/vision/models
# ADD ./tunnel/ /go/src/github.com/bigroom/vision/tunnel

WORKDIR /go/src/github.com/bigroom/vision

CMD gin -i -a=8080 -b="vision"
