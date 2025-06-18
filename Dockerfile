FROM golang:1.24-alpine AS build

ENV GO111MODULE=on

WORKDIR /src
RUN mkdir /src/bin

COPY . ./

ENV CGO_ENABLED=0
ENV GOOS=linux

RUN go build -o /src/bin/billsplitter cmd/main.go

FROM alpine:3.12.3 AS application

RUN apk --no-cache add ca-certificates

WORKDIR /opt/app
COPY --from=build /src/bin/billsplitter /opt/app/billsplitter
COPY ./internal/db /opt/app/db
COPY ./config.yml /opt/app/config.yml
