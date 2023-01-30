# docker run --rm --platform=linux/amd64 \
#   -v "$(pwd):/myapp" \
#   -w /myapp \
#   -e GOOS=linux -e GOARCH=arm \
#   golang:1.19 go build -v ./examples/simplecam

CGO_ENABLED=1 \
    GOOS=linux \
    GOARCH=arm \
    CC="zig cc -target arm-linux-musleabihf"\
    CXX="zig c++ -target arm-linux-musleabihf" go build .