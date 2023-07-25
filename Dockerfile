# Build stage
FROM golang:1.20.5-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o wss main.go


# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/wss .
# env 파일 있으면
# COPY app.env .

EXPOSE 3000
CMD [ "/app/wss" ]