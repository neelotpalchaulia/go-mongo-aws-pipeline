# Start from the official Go image
FROM golang:1.20-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the Go mod and sum files into the container
COPY go.mod go.sum ./

# Download the dependencies (this creates the go.sum file if missing)
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the Go app
RUN go build -o go-app

# Expose port 8080
EXPOSE 8080

# Run the Go application
CMD ["./go-app"]
