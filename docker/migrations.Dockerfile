FROM golang:1.25.5-trixie

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /bankapp_migrations ./cmd/migrations
RUN chmod +x /bankapp_migrations

COPY ./config/config.yaml /config.yaml
COPY ./internal/ports/db/migrations /migrations

ENTRYPOINT [ "/bankapp_migrations", "up", "-c", "/config.yaml", "-m", "/migrations" ]