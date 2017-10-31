#!/bin/bash
echo "Building Open DDNS Server..."
echo "Switching into ./web-ui to build ember project..."
cd web-ui
./node_modules/ember-cli/bin/ember build -prod

cd ..

echo "Building openddnsd-darwin..."
env GOOS=darwin GOARCH=amd64 go build -o bin/openddnsd-darwin src/*

echo "Building openddnsd-linux..."
env GOOS=linux GOARCH=amd64 go build -o bin/openddnsd-linux src/*

echo "Done!"
