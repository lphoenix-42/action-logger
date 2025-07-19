FROM golang:1.23.1-alpine AS builder

COPY . /src/

WORKDIR /src/

RUN go mod download
RUN go build -o ./bin/action_logger cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /src/bin/action_logger .

CMD ["./action"]