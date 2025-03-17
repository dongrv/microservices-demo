FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service


FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/user-service .
EXPOSE 8081
# HEALTHCHECK --interval=30s --timeout=3s CMD curl -f http://localhost:8081/health || exit 1
CMD [ "./user-service" ]
