# syntax=docker/dockerfile:1
##
## Build
##
FROM golang:1.18 AS build
LABEL project = "forum"
LABEL authors = "Willem Kuningas & Samuel Uzoagba"
LABEL version = "1.0"
# Dependancies
WORKDIR /go/src/forum
COPY go.mod ./
COPY go.sum ./
COPY . .
RUN go mod download
# Copy source files
COPY . ./
# Build
RUN go build -o ./
RUN apt-get update && apt-get install -y xdg-utils
##
## Deploy
##
EXPOSE 8080
ENTRYPOINT ["./forum"]