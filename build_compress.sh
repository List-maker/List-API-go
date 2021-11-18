#!/bin/bash -eu

VERSION="0.1.0"
OUTPUT_FILENAME="ListsBack-$VERSION"
OUTPUT_DIR="compiled"

echo "Building $OUTPUT_FILENAME ..."
go build -o $OUTPUT_DIR/$OUTPUT_FILENAME ./src && strip $OUTPUT_DIR/$OUTPUT_FILENAME && xz $OUTPUT_DIR/$OUTPUT_FILENAME
echo "Done !"