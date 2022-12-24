# Build stage
FROM golang:1.19-alpine AS builder
WORKDIR /go/src/app
RUN apk add --no-cache git

# Setup dependencies
COPY go.mod go.sum /go/src/app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go mod download
COPY ./ /go/src/app/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -v -o /go/src/app/gittlz

# Run stage
FROM alpine:3.17.0

# Install git
RUN apk add --no-cache git
RUN apk add --no-cache git-daemon

COPY --from=builder /go/src/app/gittlz /usr/bin/gittlz
LABEL Name=gittlz

EXPOSE 22
EXPOSE 80
EXPOSE 9418

# Repo directory
VOLUME /srv/git

CMD ["/usr/bin/gittlz", "serve"]