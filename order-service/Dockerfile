FROM golang:1.21-alpine
WORKDIR /app
COPY . .
RUN go mod download
COPY . .
RUN go build -o user-service
CMD [ "./order-service" ]