FROM ubuntu:latest
LABEL authors="Kibarashka"

FROM --platform=linux/arm64 golang:1.24.3-bullseye AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .


RUN CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -o /go/bin/server ./api


FROM --platform=linux/arm64 debian:bookworm-slim

RUN apt-get update && apt-get install -y ca-certificates libsqlite3-0 \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app
RUN mkdir -p data
COPY --from=builder /go/bin/server /app/server
RUN chmod +x /app/server

EXPOSE 8080
CMD ["/app/server"]
ENTRYPOINT ["top", "-b"]