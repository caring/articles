#!/bin/bash
if [ -f "$(command -v protoc)" ]; then
    VER=$(protoc --version)
    PBDIR="articles/pb/"
    echo "Using protoc version: $VER"
    protoc \
      --proto_path=$PBDIR \
      --go_out=plugins=grpc:$PBDIR \
      --go_opt=paths=source_relative $PBDIR*.proto
else
    echo "Error: protoc was not found. Please check that it is installed."
fi