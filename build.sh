#!/usr/bin/env bash
goos='linux'
goarch='amd64'

CGO_ENABLED=0 GOOS="${goos}" GOARCH="${goarch}" \
    go build -o bin/dbdef-${goos}-${goarch}
