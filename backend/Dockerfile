FROM golang:1.21-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod ./

COPY . .

RUN go mod tidy

RUN go test -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o cinematicket .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/cinematicket .

EXPOSE 8080

CMD ["./cinematicket"]
