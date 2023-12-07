FROM golang:1.21.0-alpine3.18

WORKDIR /app

COPY . ./

RUN go build -o app-cookie ./cmd

EXPOSE 8081

CMD ["./app-cookie"]