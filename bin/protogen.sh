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

script_dir=$(dirname "${BASH_SOURCE[0]}")
project_root=$(dirname "$script_dir")

# This gets the absolute path instead of relative
project_root="$( cd "$project_root" ; pwd -P )"

# protoc does not properly work with go modules
# and assumes that the project will be in your
# GOPATH. It does not have to actually be in your
# GOPATH, the directory structure just has to match.
# Check that the root directory ends with the right
# pattern so imports work as expected.
import_path="github.com/jakewright/home-automation"
if [[ ! $project_root == *"$import_path" ]]; then
    printf "\033[41mProject needs to be in a path ending with $import_path\033[0m\n"
    exit 1
fi

# Get the "GOPATH". As mentioned above, it does not
# have to be the _actual_ GOPATH. It just has to be
# the directory from which protoc can follow the
# import path to find the project's files.
# For information about how this works, see:
# https://www.tldp.org/LDP/abs/html/parameter-substitution.html
go_path=${project_root%$import_path}

# Remove slash from beginning of service directory
service_dir=${1#/}

# Create the full path to the service
service_dir="$project_root/$service_dir"

# Remove slash from the end.
service_dir=${service_dir%/}

# Make sure what we have ended up with is a directory
if [[ ! -d "$service_dir" ]]; then
    printf "\033[41mCannot find directory at $service_dir\033[0m"
    exit 1
fi

TICK="\xE2\x9C\x94"
GREEN="\033[32m"
RESET="\033[0m"

# Remove old protobuf files
find "$service_dir" -maxdepth 3 -type f -name "*.pb.go" -exec rm {} \;

# Generate all of the proto files we find
files=$(find "$service_dir" -maxdepth 3 -type f -name "*.proto")
for f in $files; do
    # Strip the project_root from the filename
    f_pretty=${f#$project_root}
    printf "Compiling $f_pretty ";

    # Generate the files
    protoc --proto_path="$go_path" --go_out="$go_path" --jrpc_out="$go_path" $f

    # Remove the ".proto" from the end of the filename
    f_base=${f%".proto"}

    printf "$GREEN$TICK$RESET\n"
done

# Run goimports on all of the generated files
find "$service_dir" -maxdepth 3 -type f -name "*.pb.go" -exec goimports -w {} \;
