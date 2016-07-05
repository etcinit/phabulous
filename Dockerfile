FROM golang
MAINTAINER Eduardo Trujillo <ed@chromabits.com>

ENTRYPOINT ["/go/bin/phabulous"]
CMD ["serve"]

EXPOSE 8085

RUN mkdir -p /go/src/github.com/etcinit/phabulous
WORKDIR /go/src/github.com/etcinit/phabulous

COPY app ./app
COPY cmd ./cmd
COPY config ./config
COPY LICENSE .

RUN go get -v -d github.com/etcinit/phabulous/cmd/phabulous \
  && go install github.com/etcinit/phabulous/cmd/phabulous
