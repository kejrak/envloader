FROM --platform=$BUILDPLATFORM golang:alpine

LABEL maintainer "kejrak"

LABEL description "Command line tool for encrypting, decrypting and loading environment variables from .ini file to binary application."

WORKDIR /app

COPY . /app

RUN apk add --no-cache make

RUN make build

ENTRYPOINT [ "./bin/envloader"]
