FROM golang:alpine

LABEL maintainer="Douglass Kirkley <doug.kirkley@gmail.com"

WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod download

COPY . .

# Build the Go app
# Build the binary.
RUN GOOS=linux GOARCH=amd64 go install