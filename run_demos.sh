#!/usr/bin/env bash

set -e

mkdir demo/output
./spiro demo/example1 demo/example1.yml demo/output
find demo/output
rm -rfv demo/output

mkdir demo/output
./spiro demo/example2 demo/example2.yaml demo/output
find demo/output
rm -rfv demo/output

mkdir demo/output
echo "x: 1" | ./spiro demo/example3 - demo/output
find demo/output
rm -rfv demo/output
