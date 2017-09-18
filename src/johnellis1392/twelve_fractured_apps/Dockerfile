# ######### #
# Version 1 #
#
# Using scratch here to build a minimal base image
# FROM scratch
# MAINTAINER John Ellis <john.ellis1392@gmail.com>
#
# Copy built binary into image
# COPY app /app
# ENTRYPOINT ["/app"]


# ######### #
# Version 2 #
# FROM alpine:latest
# MAINTAINER John Ellis <john.ellis1392@gmail.com>
# COPY app /app
# COPY entrypoint.sh /entrypoint.sh
# ENTRYPOINT ["/entrypoint.sh"]


# ######### #
# Version 3 #
#
# Multi-Stage Build:
# Build
FROM go:alpine AS build
WORKDIR /usr/src/app/src/johnellis1392/twelve_fractured_apps
ENV GOPATH=/usr/src/app

COPY . .
RUN go build .


# Run
FROM scratch
ENV GOPATH=/usr/src/app

# Configuration Variables
ENV APP_DATADIR="/var/lib/data"
ENV APP_HOST="127.0.0.1"
ENV APP_PORT="3306"
ENV APP_USERNAME="user"
ENV APP_PASSWORD="password"
ENV APP_DATABASE="test"

COPY --from=build /usr/src/app/src/johnellis1392/twelve_fractured_apps/app /usr/src/app
CMD [ "/usr/src/app" ]
