FROM ubuntu:latest

ENV GOPATH="/go"
ENV TF_TYPE="cpu"
ENV TARGET_DIRECTORY="/usr/local"

# Install Dependencies
RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y git curl build-essential

RUN mkdir -p /go/bin /go/pkg /go/src
RUN apt-get install -y golang-go


# Install Tensorflow
RUN curl -L "https://storage.googleapis.com/tensorflow/libtensorflow/libtensorflow-${TF_TYPE}-$(go env GOOS)-x86_64-1.8.0.tar.gz" \
    | tar -C ${TARGET_DIRECTORY} -xz
RUN ldconfig

RUN go get github.com/tensorflow/tensorflow/tensorflow/go
RUN go test github.com/tensorflow/tensorflow/tensorflow/go


# Build App
WORKDIR /go/src/app
COPY . .
RUN go build -o main app

ENTRYPOINT [ "./main" ]
