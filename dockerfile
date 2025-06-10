# Stage 1: Build frontend assets
FROM node:20 AS frontend
WORKDIR /frontend
COPY src/frontend/package*.json ./src/frontend/
WORKDIR /frontend/src/frontend
RUN npm ci
COPY src/frontend .
RUN npm run build

# Stage 2: Build Go binary
FROM golang:1.24 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY --from=frontend /frontend/src/frontend/dist ./src/frontend/dist
RUN go build -o luma ./src/backend

# Stage 3: Runtime image
FROM debian:bookworm-slim
WORKDIR /app
COPY --from=builder /app/luma ./luma
EXPOSE 3000
CMD ["./luma"]