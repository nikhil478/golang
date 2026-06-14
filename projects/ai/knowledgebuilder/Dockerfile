FROM golang:1.25.0 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o /run-knowledgebuilder ./cmd/main.go


FROM debian:bookworm-slim

WORKDIR /

RUN apt-get update && apt-get install -y \
    sqlite3 \
    libsqlite3-0 \
    ca-certificates

COPY --from=builder /run-knowledgebuilder /run-knowledgebuilder

CMD ["/run-knowledgebuilder", "SCRUM-7"]