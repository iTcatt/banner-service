FROM golang:1.22.2-alpine

WORKDIR /usr/src/app

RUN apk --no-cache add bash make

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY ./ ./
RUN make build
