FROM golang:1.22-alpine3.20 AS builder
WORKDIR /opt/srv/api
ENV APP_ENV=development
ENV GIT_TERMINAL_PROMPT=0
COPY . .

RUN apk update \
    && apk upgrade \
    && apk add --no-cache git mailcap

RUN go install github.com/air-verse/air@latest \
    && go mod download

CMD ["go", "run", "cmd/api/main.go"]

