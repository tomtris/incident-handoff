FROM node:22-alpine3.23 AS frontend
WORKDIR /app
COPY frontend-vue/package*.json ./
RUN npm ci
COPY frontend-vue/ .
RUN npm run build

FROM golang:1.26.4-alpine3.24 AS backend
WORKDIR /app
COPY backend-go/go.mod ./
COPY backend-go/go.sum ./
RUN go mod download
COPY backend-go/ .
COPY --from=frontend /app/dist ./frontend-vue/dist
RUN go build

FROM alpine:3.24.1
WORKDIR /app
COPY --from=backend /app/handoff ./
CMD ["./handoff"]