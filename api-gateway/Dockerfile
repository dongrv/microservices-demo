FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o api-gateway

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/api-gateway .
EXPOSE 8080
CMD [ "./api-gateway" ]
