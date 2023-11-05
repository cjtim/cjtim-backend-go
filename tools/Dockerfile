FROM golang:1.17.6-alpine as builder
WORKDIR $GOPATH/src/github.com/cjtim/cjtim-backend-go

COPY go.mod go.mod
COPY go.sum go.sum

RUN go mod tidy

COPY . .

ARG GOARCH
ARG GOOS=linux
ARG CGO_ENABLED=0
# Build the binary.
RUN go build -o main cmd/cjtim-backend-go/main.go

FROM alpine:3.15  
RUN apk update && \
    apk add --no-cache \
    ca-certificates \
    libc6-compat
WORKDIR /root/
COPY --from=builder /go/src/github.com/cjtim/cjtim-backend-go .
EXPOSE 8080
CMD ["/root/main"]