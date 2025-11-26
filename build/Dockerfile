FROM golang:1.25.3

WORKDIR /app

COPY go.mod go.sum /app/

RUN go mod download

COPY . . 

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /previewer ./cmd/server/main.go

CMD ["/app/previewer"]