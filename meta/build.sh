#!/bin/bash

# script taken from https://gist.github.com/mshafiee/5a681bbefda8f26f1f257d62f5e4a699

BIN_FILE_NAME_PREFIX=$1
PROJECT_DIR=$2
PLATFORMS="linux/386 linux/amd64 linux/arm linux/arm64 \
	   darwin/386 darwin/amd64 \
	   freebsd/386 freebsd/amd64 freebsd/arm \
	   netbsd/386 netbsd/amd64 netbsd/arm \
	   openbsd/386 openbsd/amd64 openbsd/arm \
	   windows/386 windows/amd64"

for PLATFORM in $PLATFORMS; do
        GOOS=${PLATFORM%/*}
        GOARCH=${PLATFORM#*/}
        FILEPATH="$PROJECT_DIR/artifacts/${GOOS}-${GOARCH}"
        #echo $FILEPATH
        mkdir -p $FILEPATH
        BIN_FILE_NAME="$FILEPATH/${BIN_FILE_NAME_PREFIX}"
        #echo $BIN_FILE_NAME
        if [[ "${GOOS}" == "windows" ]]; then BIN_FILE_NAME="${BIN_FILE_NAME}.exe"; fi
        CMD="CGO_ENABLED=0 GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags '-w -s' -o ${BIN_FILE_NAME}"
        #echo $CMD
        echo "${CMD}"
        eval $CMD || FAILURES="${FAILURES} ${PLATFORM}"
done
