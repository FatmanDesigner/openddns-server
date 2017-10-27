#!/bin/bash
echo "Building Open DDNS Server..."
go build -o bin/openddnsd src/*

echo "Done!"
