#!/usr/bin/env bash

cwd=`pwd`
# shellcheck disable=SC2046
# shellcheck disable=SC2164
cd $(dirname "$0")

mkdir lib
wget -c https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui.css -O lib/ui.css
wget -c https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-bundle.js -O lib/swagger.js
wget -c https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-standalone-preset.js -O lib/preset.js
wget -c https://rebilly.github.io/ReDoc/releases/latest/redoc.min.js -O lib/redoc.js

# shellcheck disable=SC2164
cd "$cwd"