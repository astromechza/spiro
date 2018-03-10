#!/usr/bin/env bash

set -e

mkdir -p build/

# add the git commit id and date
VERSION="$(git describe --tags --dirty) (commit $(git rev-parse --short HEAD) @ $(git log -1 --date=short --pretty=format:%cd))"

function buildbinary {
    export GOOS=$1
    export GOARCH=$2

    echo "Building official ${GOOS} ${GOARCH} binary for version '${VERSION}'"

    go build -i -v -o "build/spiro-${GOOS}-${GOARCH}" -ldflags "-X \"main.Version=${VERSION}\""

    echo "Done"
    ls -l "build/spiro-${GOOS}-${GOARCH}"
    file "build/spiro-${GOOS}-${GOARCH}"
    echo

    unset GOOS
    unset GOARCH
}

# platform builds
buildbinary darwin amd64
buildbinary linux amd64
buildbinary windows amd64

echo "Building local binary"
go build -i -v -ldflags "-X \"main.Version=${VERSION}\""
