#!/bin/bash

release="2018d"
file="timezones.geojson.zip"

wget "https://github.com/evansiroky/timezone-boundary-builder/releases/download/${release}/${file}"
unzip "${file}"

mkdir "reduced"
mapshaper -i dist/combined.json -simplify visvalingam 20% -clean -o reduced/reduced.json

go-bindata -pkg gotz -o tzshapefile.go reduced/

rm "${file}"
rm "dist/combined.json" "reduced/reduced.json"
rmdir "dist" "reduced"