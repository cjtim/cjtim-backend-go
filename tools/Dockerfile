FROM golang:1.17-alpine as builder
WORKDIR $GOPATH/src/github.com/cjtim/cjtim-backend-go
COPY go.mod go.mod
COPY go.sum go.sum

# Using go get.
# RUN go mod download

# Using go mod.
RUN go mod download
RUN go mod verify

# Build the binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o main cmd/cjtim-backend-go/main.go

FROM alpine:latest  
RUN apk update && \
    apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/cjtim/cjtim-backend-go .
EXPOSE 8080
CMD ["/root/main"]