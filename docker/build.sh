#!/bin/bash

# Build docker image with tag "forum"
docker image build -t forum .

# Check the exit code of the previous command
if [ $? -ne 0 ]; then
    echo "Image build failed, try again"
    exit 1
fi