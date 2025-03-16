FROM golang:1.23-alpine
WORKDIR /app
COPY . .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service
CMD [ "./user-service" ]
