#!/usr/bin/env bash

#set -xe

if [ "$1" ]; then

  # shellcheck disable=SC2046
  # shellcheck disable=SC2002
  find papaya -name '*.go' -type f -print0 | while IFS= read -r -d '' file
  do

    echo "$file"
    cat $file | grep "$1"
  done
fi