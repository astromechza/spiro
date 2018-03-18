#!/usr/bin/env bash

set -e

rm -rfv demos/output/0
rm -rfv demos/output/1
rm -rfv demos/output/2
rm -rfv demos/output/3

./spiro demos/0 demos/0/spec.yaml demos/output
find demos/output

./spiro demos/1 demos/1/spec.yaml demos/output
find demos/output

./spiro demos/2 demos/2/spec.yaml demos/output
find demos/output

echo "x: 1" | ./spiro demos/3 - demos/output
find demos/output
