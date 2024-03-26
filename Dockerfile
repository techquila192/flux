FROM golang:alpine AS builder

WORKDIR /flux_build

COPY . .

RUN go mod download

RUN go build -o flux_bin .

FROM redis:6.2-alpine AS redis

WORKDIR /flux

COPY ./config/config.json ./config/

COPY --from=builder /flux_build/flux_bin .

EXPOSE 6379

EXPOSE 5000

CMD ["sh", "-c", "redis-server --daemonize yes && ./flux_bin "]