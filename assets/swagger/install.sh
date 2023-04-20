#!/usr/bin/env bash

cwd=`pwd`
cd $(dirname "$0")

mkdir lib
wget -c https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui.css -O lib/styles.css
wget -c https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-bundle.js -O lib/scripts.js
# wget -c https://unpkg.com/swagger-ui-dist@4.5.0/swagger-ui-standalone-preset.js -O lib/preset.js
wget -c https://rebilly.github.io/ReDoc/releases/latest/redoc.min.js -O lib/redoc.js

cat lib/styles.css | base64 > lib/styles.css.b64
cat lib/scripts.js | base64 > lib/scripts.js.b64
cat lib/redoc.js | base64 > lib/redoc.js.b64

cd "$cwd"