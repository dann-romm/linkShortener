FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /app

EXPOSE 8080

CMD ["go", "run", "cmd/linkShortener/main.go"]
