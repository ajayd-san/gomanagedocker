#!/bin/bash

RELEASE=$1

if [[ -z $1 ]];
then
    echo "No release version passed"
    exit 1
else 
    echo "Building release $RELEASE"
fi 

if [ -d ./releases/ ]; then
    rm -f ./releases/*
fi 

#linux
env GOOS=linux GOARCH=amd64 go build -o ./releases/"gomanagedocker_linux_$RELEASE" github.com/ajayd-san/gomanagedocker

# macos
env GOOS=darwin GOARCH=amd64 go build -o ./releases/"gomanagedocker_darwin64_$RELEASE" github.com/ajayd-san/gomanagedocker
env GOOS=darwin GOARCH=arm64 go build -o ./releases/"gomanagedocker_darwin_arm64_$RELEASE" github.com/ajayd-san/gomanagedocker

#windows
env GOOS=windows GOARCH=amd64 go build -o ./releases/"gomanagedocker_windows64_$RELEASE.exe" github.com/ajayd-san/gomanagedocker

