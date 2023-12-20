FROM golang:alpine

WORKDIR /app

COPY . /app

RUN apk add --no-cache make

RUN make build

ENTRYPOINT [ "./bin/envLoader"]