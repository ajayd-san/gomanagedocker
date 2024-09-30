#!/bin/bash

shopt -s extglob

RELEASE=$1
if [[ -z $1 ]];
then
    echo "No release version passed"
    exit 1
else 
    echo "Building release $RELEASE"
fi 

if [ -d ./releases/ ]; then
    rm -fr ./releases/*
fi 

# linux
echo "Building linux:amd64"
env GOOS=linux GOARCH=amd64 go build -o ./releases/linux_amd64_$RELEASE/"gmd" github.com/ajayd-san/gomanagedocker

echo "Building linux:arm64"
env GOOS=linux GOARCH=arm64 go build -o ./releases/linux_arm64_$RELEASE/"gmd" github.com/ajayd-san/gomanagedocker

# MacOS
echo "Building darwin:amd64"
env GOOS=darwin GOARCH=amd64 go build -o ./releases/darwin_amd64_$RELEASE/"gmd" github.com/ajayd-san/gomanagedocker

echo "Building darwin:arm64"
env GOOS=darwin GOARCH=arm64 go build -o ./releases/darwin_arm64_$RELEASE/"gmd" github.com/ajayd-san/gomanagedocker


# Windows
echo "Building windows:amd64"
env GOOS=windows GOARCH=amd64 go build -o ./releases/windows_amd64_$RELEASE/"gmd.exe" github.com/ajayd-san/gomanagedocker

cd releases

## make a tar ball to save space
tar czf gomanagedocker_linux_amd64_$RELEASE.tar.gz linux_amd64_$RELEASE
tar czf gomanagedocker_linux_arm64_$RELEASE.tar.gz linux_arm64_$RELEASE
tar czf gomanagedocker_darwin_amd64_$RELEASE.tar.gz darwin_amd64_$RELEASE
tar czf gomanagedocker_darwin_arm64_$RELEASE.tar.gz darwin_arm64_$RELEASE
tar czf gomanagedocker_windows_amd64_$RELEASE.tar.gz windows_amd64_$RELEASE

## remove all files in releases that are not tar balls
rm -r !(*.tar.gz)
