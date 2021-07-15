#!/bin/bash

# Build the docker image and publish internal port 80 to host port 8080
# The site will be accessible on port 8080
docker build -t livestream-api . && docker run --rm -p 8080:8080 -it livestream-api
