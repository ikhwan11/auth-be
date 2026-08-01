FROM golang:1.25-alpine

WORKDIR /app

# Install git (dibutuhkan jika ada dependency dari GitHub)
RUN apk add --no-cache git

# Copy dependency terlebih dahulu agar cache Docker optimal
COPY go.mod go.sum ./

RUN go mod download

# Copy source code
COPY . .

# Build aplikasi dari cmd/api
RUN go build -o auth-be ./cmd/api

EXPOSE 8060

CMD ["./auth-be"]