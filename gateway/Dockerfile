FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .


RUN go build -o gateway-app ./gateway/main.go

FROM alpine:latest
RUN apk add --no-cache curl

WORKDIR /root/

COPY --from=builder /app/gateway-app .

CMD ["./gateway-app"]