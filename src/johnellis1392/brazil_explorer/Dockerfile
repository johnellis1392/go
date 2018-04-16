FROM golang:alpine AS build

WORKDIR /go/src/johnellis1392/brazil_explorer
RUN apk add --update git curl
RUN go get github.com/gorilla/mux

COPY . .
RUN go build -o main


FROM alpine:latest AS run

WORKDIR /usr/src/app
COPY --from=build /go/src/johnellis1392/brazil_explorer/main main

ENTRYPOINT [ "./main" ]
