#!/bin/bash

mkdir -p build

echo "::Building Linux"
export GOOS=linux
mkdir -p photos-linux
go build -o photos-linux/photos -i . &&\
tar -cvzf build/photos-linux.tar.gz photos-linux &&\
rm -r photos-linux &&\
echo Done!

echo "::Building Windows"
export GOOS=windows
mkdir -p photos-windows
go build -o photos-windows/photos.exe -i . &&\
zip build/photos-windows.zip -r photos-windows &&\
rm -r photos-windows &&\
echo Done!

echo "::Building OSX"
export GOOS=darwin
mkdir -p photos-osx
go build -o photos-osx/photos -i . &&\
tar -cvzf build/photos-osx.tar.gz photos-osx &&\
rm -r photos-osx &&\
echo Done!
