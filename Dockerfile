# Build stage
FROM golang:1.21-alpine3.18 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .
COPY wait-for.sh .
COPY static/default_avatar.png ./static
COPY db/migration ./db/migration

EXPOSE 8080
CMD [ "/app/main" ]

# entrypoint 會把 cmd 的指令帶入其中
ENTRYPOINT [ "/app/start.sh" ]