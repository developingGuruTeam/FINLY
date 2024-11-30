FROM golang:1.23.2 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /cashManagerApp app/cmd/main.go

CMD ["sh", "-c", "sleep 10 && /cashManagerApp"]
