# Build stage
FROM golang:1.22.4 AS builder

# Set the Current Working Directory inside the container
WORKDIR /HoneyMesh

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies. Dependencies will be cached if the go.mod and go.sum files are not changed
RUN go mod tidy
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Build the Go app
RUN go build -o main .

COPY ./assets /HoneyMesh/assets
COPY ./templates /HoneyMesh/templates

# Use Ubuntu as base image for the final stage
FROM ubuntu:22.04

# Install necessary tools
RUN apt-get update && \
    apt-get install -y wget gnupg2 lsb-release netcat bash

# Add MongoDB repository
RUN wget -qO - https://www.mongodb.org/static/pgp/server-6.0.asc | apt-key add - && \
    echo "deb [ arch=amd64,arm64 ] https://repo.mongodb.org/apt/ubuntu $(lsb_release -cs)/mongodb-org/6.0 multiverse" | tee /etc/apt/sources.list.d/mongodb-org-6.0.list

# Install MongoDB
RUN apt-get update && \
    apt-get install -y mongodb-org

# Set the Current Working Directory inside the container
WORKDIR /HoneyMesh

# Copy the Pre-built binary file from the previous stage
COPY --from=builder /HoneyMesh/main .

COPY ./assets /HoneyMesh/assets
COPY ./templates /HoneyMesh/templates

# Create directory for MongoDB data
RUN mkdir -p /data/db

# Expose ports 8080 and 27017
EXPOSE 8080 27017

# Copy and run start script
COPY start.sh /start.sh
RUN chmod +x /start.sh

# Command to run the start script
CMD ["/start.sh"]