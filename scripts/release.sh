#!/bin/bash

if [[ -z "$TAG_NAME" ]]; then
    TAG_NAME=$(git rev-parse --short HEAD || echo "Unknown")
fi

export CGO_ENABLED=0
platforms=("linux" "android")
architectures=("amd64" "386" "arm" "arm64")

rm -rf ./build/releases
mkdir -p ./build/releases

for os in "${platforms[@]}"; do

    for arch in "${architectures[@]}"; do
        output_name="anping-${os}-${arch}"
        if [ "$os" == "windows" ]; then
            output_name="${output_name}.exe"
        fi
        # if [ "$os" == "android" ]; then
        #     export CGO_ENABLED=1
        # else
        #     export CGO_ENABLED=0
        # fi

        echo "Building ${output_name}..."
        GOOS="$os" GOARCH="$arch" go build -o "build/releases/${output_name}" \
            -tags "" \
            -trimpath -ldflags "-w -s -X main.version=${TAG_NAME} -buildid=" ./cmd/anping
    done

done
