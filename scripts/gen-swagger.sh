#!/bin/bash
set -e

SOURCE="${BASH_SOURCE[0]}"
ROOT_DIR="$(cd -P "$(dirname "$SOURCE")/.." && pwd)"

swag init \
    --dir $ROOT_DIR/cmd/api \
    --output $ROOT_DIR/docs \
    --parseDependency \
    --parseInternal \
    --parseDepth 2
