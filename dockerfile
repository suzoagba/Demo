# Use a larger base image for the build stage
FROM golang:1.18 AS build
LABEL project="forum"
LABEL authors=""
LABEL version="2.0"

# Install necessary packages for building C-based dependencies
RUN apt-get update && apt-get install -y build-essential

# Install xdg-utils package
RUN apt-get install -y xdg-utils

# Set the working directory
WORKDIR /go/src/forum

# Copy go.mod and go.sum files separately to leverage caching
COPY go.mod .
COPY go.sum .
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -o forum .

# Use a smaller base image for the deployment stage
FROM alpine:latest AS deploy
LABEL project="forum"
LABEL authors=""
LABEL version="2.0"

# Install required packages
RUN apk --no-cache add ca-certificates xdg-utils

# Copy the binary from the build stage
COPY --from=build /go/src/forum/forum .

# Expose the necessary port
EXPOSE 8080

# Set the entrypoint command
ENTRYPOINT ["./forum"]