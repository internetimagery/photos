#!/bin/bash

# Make publish directory
mkdir -p publish

echo "::Building Linux"
export GOOS=linux
go build -v -o publish/photos-linux

echo "::Building Windows"
export GOOS=windows
go build -v -o publish/photos-windows.exe

echo "::Building OSX"
export GOOS=darwin
go build -v -o publish/photos-osx
