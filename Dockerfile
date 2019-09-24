FROM golang:alpine as build-env
MAINTAINER Chen Lei "my@mysq.to"

RUN mkdir /log
WORKDIR /log
# <- COPY go.mod and go.sum files to the workspace
COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download
# COPY the source code as the last step
COPY . .
RUN go build -o /go/bin/example example/example.go

# to build minimal image
FROM alpine:latest
RUN mkdir /app
WORKDIR /app
COPY --from=build-env /go/bin/example /app/example

# if using golang:latest to build, uncomment the following line to fix 64 bit binary "not found" error.
# RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
ENTRYPOINT ["/app/example"]
