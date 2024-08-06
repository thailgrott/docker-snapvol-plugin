FROM golang:1.18 as builder

WORKDIR /app

COPY . .

RUN go build -o snapvol main.go plugin_api.go app/btrfs_manager.go

FROM alpine:latest

WORKDIR /usr/local/bin

COPY --from=builder /app/snapvol .

ENTRYPOINT ["./snapvol"]
