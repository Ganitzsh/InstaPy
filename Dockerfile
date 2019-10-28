FROM golang:1.13 as build

WORKDIR /app

COPY . .

RUN go mod download
RUN go build -tags netgo -a -v -o /bin/gobot

FROM alpine as runtime

WORKDIR /app

COPY --from=build /bin/gobot /bin/gobot
COPY --from=build /app .

RUN chmod +x /bin/gobot
RUN echo "{\"username\": \"user\",\"password\":\"password\"}" > config.json
RUN env | grep PATH && ls -l | grep config.json && cat config.json

EXPOSE 8080

ENTRYPOINT ["/bin/gobot"]
