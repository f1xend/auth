FROM golang:1.21.5-alpine AS builder

COPY . /github.com/f1xend/auth/source/
WORKDIR /github.com/f1xend/auth/source/

RUN go mod download
RUN go build -o ./bin/auth_server cmd/grpc_server/main.go

FROM alpine:latest

WORKDIR /root/
COPY --from=builder /github.com/f1xend/auth/source/bin/auth_server .

CMD ["./auth_server"]