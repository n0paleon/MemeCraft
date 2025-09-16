# stage 1: build the binary
FROM golang:1.25-alpine AS builder
WORKDIR /app
# install build dependencies
RUN apk add --no-cache git
# copy all source code and project files to working directory
COPY . .
# download required go modules
RUN go mod download
# build the binary
RUN mkdir -p bin
RUN go build -ldflags="-s -w" -o bin/memecraft main.go

# stage 2: runtime image
FROM alpine:latest
WORKDIR /app
# copy the binary and project files from the build stage
COPY --from=builder /app .
# set the entrypoint and command
ENTRYPOINT ["bin/memecraft"]