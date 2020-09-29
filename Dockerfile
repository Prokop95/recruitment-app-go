FROM golang:1.14.4 AS builder

ENV GO111MODULE=on \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /build
RUN mkdir -p /build/config

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY /config/app.yaml /config/
COPY . .

RUN go build -o recruitment-api ./cmd/main.go

WORKDIR /dist

RUN cp /build/recruitment-api .
RUN cp -r /build/config/ ./config/

FROM ubuntu:18.04

COPY --from=builder /dist/ /
CMD ["./recruitment-api"]
