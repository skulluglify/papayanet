#!/usr/bin/env bash

set -xe

cwd=`pwd`
cd $(dirname "$0")

echo 'package binaries' >assets.go

echo 'var SwagBinStyles = "data:text/css;base64," +' >>assets.go

cat lib/styles.css.b64 | while IFS= read -r x
do

  echo '"'"$x"'" +' >>assets.go
done

echo '""' >>assets.go

echo 'var SwagBinScripts = "data:application/javascript;base64," +' >>assets.go

cat lib/scripts.js.b64 | while IFS= read -r x
do

  echo '"'"$x"'" +' >>assets.go
done

echo '""' >>assets.go

# echo 'var SwagBinPresets = "data:application/javascript;base64," +' >>assets.go


# cat lib/preset.js.b64 | while IFS= read -r x
# do

#   echo '"'"$x"'" +' >>assets.go
# done

# echo '""' >>assets.go

echo 'var SwagBinRedocScripts = "data:application/javascript;base64," +' >>assets.go

cat lib/redoc.js.b64 | while IFS= read -r x
do

  echo '"'"$x"'" +' >>assets.go
done

echo '""' >>assets.go

cd "$cwd"