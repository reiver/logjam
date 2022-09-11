#!/bin/bash

echo "Starting Logjam Server"
cd ./backend && go run . -vvvvvv -http-port 8090