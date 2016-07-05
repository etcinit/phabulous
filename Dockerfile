FROM golang
RUN go get github.com/etcinit/phabulous/cmd/phabulous
ADD config /go
EXPOSE 8085
CMD phabulous serve
