#!/bin/sh
mkdir -p ./build
for f in ./plugins/*.go; do
  if [[ ! -z "$f" ]]; then
    file=$(basename $f .go)
    echo "Building $file.so"
    go build -buildmode=plugin -o ./build/$file.so ./plugins/$file.go
  fi
done
