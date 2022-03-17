FROM alpine:3.15

RUN apk add --no-cache libc6-compat
WORKDIR /root/
COPY ./main ./main

EXPOSE 8080
CMD ["/root/main"]