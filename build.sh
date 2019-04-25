#!/usr/bin/env sh

dir=$(
    cd $(dirname $0)
    pwd
)

now=$(date +%Y-%m-%dT%H:%M:%S%z)
version=$(cat VERSION)
dist=dist/fproxy

rm -rf ./dist/

docker run --rm \
    -v /Users/shu/go:/go \
    -v ${dir}:/build \
    -w /build \
    -e GOPROXY=https://goproxy.amzcs.com \
    shuxs/golang:builder \
    sh -c "go build -v -ldflags='-s -w -X main.Version=${version} -X main.BuildTime=${now}' -o ${dist} && upx -9 -o ${dist}-min ${dist} && ls -la dist"

if [[ $? == 0 ]]; then
     docker build -t shuxs/fproxy:latest . &&
     docker push shuxs/fproxy:latest
fi
