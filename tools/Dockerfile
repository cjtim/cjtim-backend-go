FROM golang:1.16 as builder
WORKDIR $GOPATH/src/github.com/cjtim/cjtim-backend-go
COPY . .

# Using go get.
RUN go get -d -v

# Using go mod.
# RUN go mod download
# RUN go mod verify
# Build the binary.
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o main cmd/cjtim-backend-go/main.go

FROM alpine:latest  
RUN apk update && \
    apk add --no-cache ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/cjtim/cjtim-backend-go .
EXPOSE 8080
CMD ["/root/main"]