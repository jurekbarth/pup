FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY ./pupauth ./pupauth
COPY ./main .
CMD ["./main"]
