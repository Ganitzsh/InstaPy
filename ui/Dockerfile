FROM golang:latest as builder

ADD . /go/src/instabot
WORKDIR /go/src/instabot
RUN CGO_ENABLED=0 GOOS=linux go build -mod=vendor -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/instabot/app .
CMD ["./app"]

