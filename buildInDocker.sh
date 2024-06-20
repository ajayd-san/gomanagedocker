#!/bin/bash

# this script is used to build release binaries inside an older version of golang docker container
# This is cuz I run a rolling release distro (I use arch btw) and run the latest version of glibc
# However, many LTS distros run older versions of glibc and the binary will not launch in those distros. 
# ref issue `https://github.com/ajayd-san/gomanagedocker/issues/8`
#


RELEASE=$1
if [[ -z $1 ]];
then
    echo "No release version passed"
    exit 1
fi 

docker run --rm -v "$PWD":/src -w /src --env GOFLAGS="-buildvcs=false" golang:1-bullseye ./buildRelease.sh $RELEASE
