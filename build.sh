#!/bin/bash

# Make publish directory
mkdir -p publish

echo "::Building Linux"
export GOOS=linux
go build -o publish/photos-linux -i . && echo Done!

echo "::Building Windows"
export GOOS=windows
go build -o publish/photos-windows.exe -i . && echo Done!

echo "::Building OSX"
export GOOS=darwin
go build -o publish/photos-osx -i . && echo Done!
