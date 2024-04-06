FROM golang:1.22.1-alpine

RUN mkdir /app
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o app ./cmd/sso/main.go

ENTRYPOINT ["./app"]
