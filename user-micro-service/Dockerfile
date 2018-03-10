FROM golang:1.10

MAINTAINER Adrien Domoison "adomoison@gmail.com"

WORKDIR /go/src/github.com/adriendomoison/apigoboot/user-micro-service

CMD go build component/serve-micro-service.go && ./serve-micro-service

EXPOSE 4200