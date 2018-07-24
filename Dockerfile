FROM golang:1.9
ENV GOPATH=/usr/go

COPY . /usr/go/src/github.com/hitesh-goel/test-news-api
WORKDIR /usr/go/src/github.com/hitesh-goel/test-news-api
RUN go get -v
RUN go install -v

CMD ["/usr/go/bin/test-news-api"]
EXPOSE 8080