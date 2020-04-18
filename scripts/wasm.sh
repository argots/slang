#!/bin/bash

set -ex

GOOS=js GOARCH=wasm go build -o ./bin/example.wasm ./cmd/example/
