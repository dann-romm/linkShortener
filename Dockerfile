FROM golang:1.18-alpine

WORKDIR /app

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . /app

EXPOSE 8080

# TODO: think about how to build app in a different ways
# Makefile maybe

#RUN go build -o weatherBackend cmd/weatherBackend/main.go
#CMD ["./weatherBackend"]
