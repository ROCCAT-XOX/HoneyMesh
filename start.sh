#!/bin/bash

# Start MongoDB in the background
mongod --fork --logpath /var/log/mongodb.log --dbpath /data/db

# Wait for MongoDB to be available
until nc -z localhost 27017; do
  echo "Waiting for MongoDB to start..."
  sleep 1
done

# Start the Go app
./main