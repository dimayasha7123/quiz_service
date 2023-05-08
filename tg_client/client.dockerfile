# v3
FROM golang:alpine as builder
RUN apk --no-cache add ca-certificates
RUN mkdir /build/
ADD . /build/
WORKDIR /build/
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o clientbin ./tg_client/cmd/client
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /build/clientbin /app/
COPY ../.env /app/
COPY ../tg_bot_config.yaml /app/
WORKDIR /app
ENTRYPOINT ["./clientbin", "--env", ".env", "--config", "tg_bot_config.yaml"]

# v2
# FROM golang:alpine as builder
# RUN mkdir /build
# ADD . /build/
# WORKDIR /build
# RUN go build -o main ./cmd/server
# FROM alpine
# COPY --from=builder /build/main /app/
# COPY /config/config.yaml /app/
# WORKDIR /app
# ENTRYPOINT ["./main", "--config", "./config.yaml"]


# v1
# FROM golang:alpine
# RUN mkdir /app
# ADD . /app/
# WORKDIR /app
# RUN go build -o main ./cmd/server
# ENTRYPOINT "./main"