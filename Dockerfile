FROM golang:latest AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/server/main.go

FROM alpine:latest

RUN apk update && apk upgrade
RUN apk --no-cache add ca-certificates

COPY --from=builder /app/main /app/main

COPY example.config.yaml ./config.yaml
COPY api/ ./api/
COPY swaggerui/ ./swaggerui/

CMD ["/app/main"]
