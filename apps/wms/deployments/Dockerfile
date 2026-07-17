# I56 WMS — Production Dockerfile
# Stage 1: Build Go binary
FROM golang:1.24-alpine AS go-builder
WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -ldflags="-s -w" -o /i56-server ./cmd/server/

# Stage 2: Build frontend
FROM node:20-alpine AS fe-builder
WORKDIR /fe
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 3: Runtime
FROM alpine:3.20
RUN apk add --no-cache ca-certificates tzdata nginx curl
ENV TZ=Asia/Shanghai
ENV APP_ENV=production
ENV APP_DEBUG=false
ENV SERVER_PORT=8080

COPY --from=go-builder /i56-server /opt/i56/i56-server
COPY --from=fe-builder /fe/dist /opt/i56/frontend

COPY deployments/nginx.conf /etc/nginx/http.d/default.conf
COPY deployments/entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh /opt/i56/i56-server

EXPOSE 80
ENTRYPOINT ["/entrypoint.sh"]
