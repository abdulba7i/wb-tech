FROM golang:1.24.4-bullseye

RUN apt-get update && apt-get install -y \
    librdkafka-dev \
    && rm -rf /var/lib/apt/lists/*

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN go build -o api ./cmd/api/main.go

COPY wait-for-it.sh ./wait-for-it.sh
RUN chmod +x ./wait-for-it.sh