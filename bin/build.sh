#!/bin/bash -eu

if [ ! -L "/bin/podcast-dl" ]; then
	ln -s /workspace/rel/podcast-dl_linux-amd64 /bin/podcast-dl
fi

source /build-common.sh

BINARY_NAME="podcast-dl"
COMPILE_IN_DIRECTORY="cmd/podcast-dl"

standardBuildProcess
