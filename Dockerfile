# Use an official Golang runtime as a base image
FROM golang:alpine

# Set the working directory inside the container
WORKDIR /app

# Copy the local source code to the container's working directory
COPY . .

# Build the Go application
RUN go build -o app

# Expose the port that the application listens on (replace 8080 with the actual port your application listens on)
EXPOSE 8080 

# Command to run the application when the container starts
CMD ["./app"]
