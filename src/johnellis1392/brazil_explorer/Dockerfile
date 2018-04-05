# ############# #
# ### Build ### #
# ############# #
FROM golang:alpine AS build

RUN apk --update add git curl
RUN go get github.com/gorilla/mux

WORKDIR /go/src/app
COPY . .
RUN go build -o main *.go


# ########### #
# ### Run ### #
# ########### #
FROM alpine:latest AS run

# Mountable Volume for Retrieving Host Workspace Information
VOLUME /workplace

WORKDIR /usr/src/app
COPY --from=build /go/src/app/main ./main

ENTRYPOINT [ "/usr/src/app/main" ]
