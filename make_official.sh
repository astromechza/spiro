#!/usr/bin/env bash

set -e

# first build the version string
VERSION_NUM=0.1

# add the git commit id and date
VERSION="$VERSION_NUM (commit $(git rev-parse --short HEAD) @ $(git log -1 --date=short --pretty=format:%cd))"

# change these!
GO_PROJECT=github.com/AstromechZA/go-cli-template
BINARY_NAME=go-cli-template

function buildbinary {
    goos=$1
    goarch=$2

    echo "Building official $goos $goarch binary for version '$VERSION'"

    outputfolder="build/${goos}_${goarch}"
    echo "Output Folder $outputfolder"
    mkdir -pv $outputfolder

    export GOOS=$goos
    export GOARCH=$goarch

    go build -i -v -o "$outputfolder/$BINARY_NAME" -ldflags "-X \"main.VersionString=$VERSION\"" $GO_PROJECT

    echo "Done"
    ls -lh "$outputfolder/$BINARY_NAME"
    file "$outputfolder/$BINARY_NAME"
    echo
}

# build local 
unset GOOS
unset GOARCH
go build -ldflags "-X \"main.VersionString=$VERSION\"" $GO_PROJECT

# build for mac
buildbinary darwin amd64

# build for linux
buildbinary linux amd64

# zip up
tar -czf $BINARY_NAME-${VERSION_NUM}.tgz -C build .
ls -lh $BINARY_NAME-${VERSION_NUM}.tgz
file $BINARY_NAME-${VERSION_NUM}.tgz