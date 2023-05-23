#!/usr/bin/env bash

set -xe

echo -n >metadata.lst
# shellcheck disable=SC2046
# shellcheck disable=SC2002
find papaya -name '*.go' -type f -print0 | while IFS= read -r -d '' file
do

  echo $(cat "$file" | sha1sum | awk '{print $1}') "$file" >>metadata.lst
done
