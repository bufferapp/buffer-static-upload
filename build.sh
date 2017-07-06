#!/bin/bash

# Clean
echo "Cleaning up"
rm -rf dist >/dev/null 2>&1
mkdir dist >/dev/null 2>&1
rm buffer-static-upload >/dev/null 2>&1

# Build native
echo "Compiling native binary"
go build main.go
mv main buffer-static-upload

# Build linux
echo "Compiling linux binary"
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" main.go
mv main ./dist/buffer-static-upload-Linux

# Build mac
echo "Compiling mac binary"
env GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" main.go
mv main ./dist/buffer-static-upload-Darwin

echo "Done!"
