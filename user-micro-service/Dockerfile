FROM golang:1.10

MAINTAINER Adrien Domoison "adomoison@gmail.com"

WORKDIR /go/src/github.com/adriendomoison/apigoboot/user-micro-service

COPY . .

RUN go get -u github.com/kardianos/govendor

RUN govendor sync

RUN GOOS=linux GOARCH=arm go build -v component/serve-micro-service.go

CMD ./serve-micro-service

EXPOSE 4200