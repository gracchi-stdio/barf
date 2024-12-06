FROM golang:1.23-alpine

RUN apk add --no-cache nodejs npm nginx

WORKDIR /app

# Go backend
COPY go.mod go.sum ./
RUN go mod download
COPY . .

# Frontend
WORKDIR /app/web
COPY web/package*.json ./
RUN npm ci
COPY web/ .
RUN npm run build

WORKDIR /app
RUN go build -o barf-api ./cmd/http

RUN cp -r /app/web/dist/* /usr/share/nginx/html/
COPY web/nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 3000 8080

CMD ["sh", "-c", "nginx && ./barf-api"]