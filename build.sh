#!/bin/bash
echo "Building Open DDNS Server..."
echo "Switching into ./web-ui to build ember project..."
cd web-ui
./node_modules/ember-cli/bin/ember build -prod

cd ..

echo "Building openddnsd..."
go build -o bin/openddnsd src/*

echo "Done!"
