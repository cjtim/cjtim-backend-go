FROM golang:1.15

WORKDIR /go/src/github.com/cjtim/cjtim-backend-go
COPY . .
RUN go build main.go

EXPOSE 8080
CMD ["./main"]

