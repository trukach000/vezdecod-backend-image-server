
FROM golang:1.15-alpine
LABEL maintainer="vezdecod"

RUN apk --no-cache add tzdata

WORKDIR /go/src

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . /go/src
WORKDIR /go/src

RUN go build -o api.server main.go

ENV DB_HOST="host.docker.internal"
ENV HOST="0.0.0.0"

ENTRYPOINT ["./api.server"]
