#!/bin/sh
while true; do
  echo "Running the application at $(date)"
  /go/main 
  echo "Sleeping for 1 hour..."
  sleep 3600
done
