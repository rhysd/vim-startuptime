#! /bin/bash

set -e

rm -rf release
gox -verbose
mkdir -p release
mv vim-startuptime_* release/
cd release
for bin in *; do
    if [[ "$bin" == *windows* ]]; then
        command="vim-startuptime.exe"
    else
        command="vim-startuptime_"
    fi
    mv "$bin" "$command"
    zip "${bin}.zip" "$command"
    rm "$command"
done
