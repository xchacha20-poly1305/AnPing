#!/bin/bash

if [[ -z "$TAG_NAME" ]]; then
    TAG_NAME=$(git rev-parse --short HEAD || echo "Unknown")
fi

export CGO_ENABLED=0
platforms=("linux" "windows" "darwin" "freebsd" "openbsd")
architectures=("amd64" "386" "arm" "arm64" "loong64")

rm -rf ./build/releases
mkdir -p ./build/releases

for os in "${platforms[@]}"; do

    for arch in "${architectures[@]}"; do
        output_name="anping-${os}-${arch}"
        if [ "$os" == "windows" ]; then
            output_name="${output_name}.exe"
        fi

        [[ "$arch" == "loong64" ]] && [[ "$os" != "linux" ]] && continue

        echo "Building ${output_name}..."
        GOOS="$os" GOARCH="$arch" go build -v -o "build/releases/${output_name}" \
            -tags "" \
            -trimpath -ldflags "-w -s -X main.version=${TAG_NAME} -buildid=" ./cmd/anping
    done

done
