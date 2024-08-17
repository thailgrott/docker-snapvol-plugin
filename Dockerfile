FROM golang:1.18 as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY app/main.go .
COPY app/plugin_api.go .
COPY app/btrfs_manager.go .

RUN CGO_ENABLED=0 go build -o snapvol main.go plugin_api.go btrfs_manager.go

FROM alpine:latest

WORKDIR /usr/local/bin

COPY config.json .
COPY --from=builder /app/snapvol .

ENTRYPOINT ["./snapvol"]
