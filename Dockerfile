FROM --platform=$BUILDPLATFORM golang:alpine

LABEL maintainer "kejrak"

WORKDIR /app

COPY . /app

RUN apk add --no-cache make

RUN make build

ENTRYPOINT [ "./bin/envLoader"]
