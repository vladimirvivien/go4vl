#! /bin/bash

./fix-mods.sh

# Cross-compile using zig to build for 32-bit arm
CGO_ENABLED=1 GOOS=linux GOARCH=arm GOARM=7 CC="zig cc -target arm-linux-musleabihf" CXX="zig c++ -target arm-linux-musleabihf" go build -o webcam .

