FROM golang:1.22.5 AS builder

WORKDIR /app

ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GO11MODULE=on

COPY . .

RUN go mod download
RUN go build -o server cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/server ./server
COPY --from=builder /app/global-bundle.pem ./global-bundle.pem

CMD [ "./server" ]