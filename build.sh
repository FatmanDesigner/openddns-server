#!/bin/bash
echo "Building Open DDNS Server..."
env GOOS=darwin GOARCH=amd64 go build -o bin/openddnsd-darwin src/*
env GOOS=linux GOARCH=amd64 go build -o bin/openddnsd-linux src/*

echo "Done!"
