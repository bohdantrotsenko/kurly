#!/bin/bash

# script taken from https://gist.github.com/mshafiee/5a681bbefda8f26f1f257d62f5e4a699



BIN_FILE_NAME_PREFIX=$1
PROJECT_DIR=$2
PLATFORMS="linux/amd64 linux/arm \
	   darwin/amd64 \
	   freebsd/amd64 freebsd/arm \
	   netbsd/amd64 netbsd/arm \
	   openbsd/amd64 openbsd/arm \
	   windows/amd64"

apt update
apt install -y zip

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
	if [[ "${GOOS}" == "windows" ]]; then
	    zip $FILEPATH/${GOOS}-${GOARCH}.zip -j ${BIN_FILE_NAME} LICENSE README.md
	    rm $FILEPATH/kurly.exe
	else
	    cp README.md LICENSE meta/kurly.man $FILEPATH
	    tar czvf $FILEPATH/${GOOS}-${GOARCH}.tar.gz -C $FILEPATH kurly README.md LICENSE kurly.man
	    rm $FILEPATH/kurly $FILEPATH/README.md $FILEPATH/LICENSE $FILEPATH/kurly.man
	fi
done
