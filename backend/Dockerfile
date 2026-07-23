# ─── Build stage ───
FROM golang:1.22-alpine AS build
WORKDIR /src

# 国内网络加速
# RUN go env -w GOPROXY=https://goproxy.cn,direct

COPY go.mod go.sum* ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /out/ticket-server .

# ─── Runtime stage ───
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata && \
    addgroup -S app && adduser -S app -G app
WORKDIR /app
COPY --from=build /out/ticket-server /app/ticket-server
RUN mkdir -p /app/data && chown -R app:app /app
USER app
EXPOSE 8000
ENV DATABASE_URL=/app/data/tickets.db
ENTRYPOINT ["/app/ticket-server"]
