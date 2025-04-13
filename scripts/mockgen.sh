#!/usr/bin/env bash

mockgen_cmd="mockgen"

find . -type f -path '*/repository/*' -name '*.go' -and -not -name '*_test.go' -and -not -path '*/testutil/*' | while read -r source_file; do
    source_dir=$(dirname "$source_file")
    destination_dir="$source_dir/testutil"

    source_filename=$(basename "$source_file")
    destination_file="$destination_dir/${source_filename%.go}_mock.go"

    package=$(basename "$source_dir")

    mkdir -p "$destination_dir"
    $mockgen_cmd -source "$source_file" -package "testutil" -destination "$destination_file"
done