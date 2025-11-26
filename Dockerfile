FROM golang:1.25.3

WORKDIR /app

COPY go.mod go.sum /app/

COPY cmd /app/cmd

COPY internal /app/internal

COPY pkg /app/pkg

RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /previewer ./cmd/server/main.go

CMD ["/previewer"]