FROM golang:1.22.1

RUN mkdir /app
WORKDIR /app
COPY . .

RUN go get -d -v ./...
RUN go install -v ./...

RUN go build -o app ./cmd/sso/main.go

ENTRYPOINT ["./app", "--config=/app/config/local.yaml"]
