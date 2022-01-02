FROM golang:1.15.7-alpine3.13

ADD . /aseri
WORKDIR /aseri/servers/httpclient

RUN apk add git
RUN go mod download
RUN go build -o main .

EXPOSE 443

CMD  ["/aseri/servers/httpclient/main"]