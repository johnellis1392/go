FROM golang:alpine AS build

WORKDIR /go/src/app
RUN apk add --update git curl
RUN go get github.com/gorilla/mux

COPY . .
RUN go build -o main *.go


FROM alpine:latest AS run

WORKDIR /usr/src/app
COPY --from=build /go/src/app/main main

ENTRYPOINT [ "./main" ]
