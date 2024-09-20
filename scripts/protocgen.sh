#!/usr/bin/env bash
set -e

echo "Generating Protocol Buffer code..."
cd proto
proto_dirs=$(find ./connect -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  for file in $(find "${dir}" -maxdepth 1 -name '*.proto'); do
    if grep go_package $file &> /dev/null ; then
      buf generate --template buf.gen.gogo.yaml $file
    fi
  done
done

cd ..

# move proto files to the right places
cp -r github.com/skip-mev/connect/v2/* ./
rm -rf github.com

# go mod tidy --compat=1.20
