# Global build arguments for all stages
ARG GOOS=linux
ARG GOARCH=amd64

FROM golang:1.16 AS builder
ARG GOOS
ARG GOARCH
ENV GOOS=$GOOS
ENV GOARCH=$GOARCH
WORKDIR /go/src/app
COPY . .
RUN make build

FROM ubuntu:20.04
ARG GOOS
ARG GOARCH
COPY --from=builder /go/src/app/fibo_${GOOS}_${GOARCH} /usr/bin/fibo