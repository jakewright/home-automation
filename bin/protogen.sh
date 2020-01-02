# Usage: bin/protogen.sh [service.foo]

set -e

# Check for protoc
if [[ ! -x "$(command -v protoc)" ]]; then
    printf "\033[41mPlease install protoc\033[0m\n"
    exit 1
fi

# Check for protoc-gen-go
if [[ ! -x "$(command -v protoc-gen-go)" ]]; then
    printf "\033[41mPlease install protoc-gen-go\033[0m\n"
    exit 1
fi

# Check for protoc-gen-jrpc
if [[ ! -x "$(command -v protoc-gen-jrpc)" ]]; then
    printf "\033[41mPlease install protoc-gen-jrpc\033[0m (run tools/install)\n"
    exit 1
fi

DIR=$(dirname "${BASH_SOURCE[0]}")
ROOT=$(dirname "$DIR")
ROOT_ABS="$( cd "$ROOT" ; pwd -P )"

# protoc does not properly work with go modules
# and assumes that the project will be in your
# GOPATH. It does not have to actually be in your
# GOPATH, the directory structure just has to match.
# Check that the root directory ends with the right
# pattern so imports work as expected.
IMPORT_PATH="github.com/jakewright/home-automation"
if [[ ! $ROOT_ABS == *"$IMPORT_PATH" ]]; then
    printf "\033[41mProject needs to be in a path ending with $IMPORT_PATH\033[0m\n"
    exit 1
fi

# Get the "GOPATH". As mentioned above, it does not
# have to be the _actual_ GOPATH. It just has to be
# the directory from which protoc can follow the
# import path to find the project's files.
# For information about how this works, see:
# https://www.tldp.org/LDP/abs/html/parameter-substitution.html
SRC=${ROOT_ABS%$IMPORT_PATH}

# Remove slash from beginning of service directory
SVC_DIR=${1#/}

# Create the full path to the service
SVC_DIR="$ROOT_ABS/$SVC_DIR"

# Remove slash from the end.
SVC_DIR=${SVC_DIR%/}

# Make sure what we have ended up with is a directory
if [[ ! -d "$SVC_DIR" ]]; then
    printf "\033[41mCannot find directory at $SVC_DIR\033[0m"
    exit 1
fi

TICK="\xE2\x9C\x94"
GREEN="\033[32m"
RESET="\033[0m"

# Generate all of the proto files we find
FILES=$(find "$SVC_DIR" -maxdepth 3 -type f -name "*.proto")
for f in $FILES; do
    PRETTY_NAME=${f#$ROOT_ABS}
    printf "Compiling $PRETTY_NAME ";
    protoc --proto_path="$SRC" --go_out="$SRC" --jrpc_out="$SRC" $f
    printf "$GREEN$TICK$RESET\n"
done
