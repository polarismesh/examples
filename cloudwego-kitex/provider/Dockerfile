############################################################
# Dockerfile to build Kitex-Polaris quickstart example provider

# First of all, you need to build the binary from the source code,
# use `CGO_ENABLED=0 go build -ldflags "-s -w" -o provider` to build the binary.

############################################################

FROM golang:alpine

WORKDIR /root
ADD . /root

RUN chmod +x ./provider

ENTRYPOINT ./provider