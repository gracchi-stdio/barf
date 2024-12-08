FROM golang:1.23-alpine

# Install Go
RUN apk add --no-cache go git npm

# Set Go environment variables
ENV GOPATH /go
ENV PATH $GOPATH/bin:$PATH

# Create necessary directories
RUN mkdir -p /app /go/bin

WORKDIR /app

# Install Air for hot reload
RUN go install github.com/air-verse/air@latest

# Copy Go files
COPY go.mod go.sum ./
RUN go mod download


# Set up frontend
WORKDIR /app/web
COPY web/package*.json ./
RUN npm install

# Copy the rest of the application
WORKDIR /app
COPY . .

# Install frontend dependencies
RUN cd /app/web && npm install

CMD ["air", "-c", ".air.toml", "npm", "run", "dev"]
