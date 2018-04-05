FROM golang:alpine AS build

ENV GOOS=darwin
ENV GOARCH=amd64

RUN apk --update add git
RUN go get github.com/lxn/walk

WORKDIR /go/src/app
COPY . .

RUN go build -o main main.go


FROM alpine:latest AS run

VOLUME /usr/src/app
WORKDIR /usr/src/app
COPY --from=build /go/src/app/main /tmp/main

CMD [ "cp", "/tmp/main", "/usr/src/app/main" ]
