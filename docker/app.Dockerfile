FROM golang:1.25.5-trixie

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o /bankapp ./cmd/bankapp
RUN chmod +x /bankapp

ENV APP_PORT=8080

EXPOSE $APP_PORT

ENTRYPOINT [ "/bankapp" ]