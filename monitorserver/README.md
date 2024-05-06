# How to use

1. download [nginx for windows](https://nginx.org/en/download.html) to a directory
1. build this project: `go build -ldflags="-s -w"`
1. copy `monitor.exe`, `template.html`, `conf/nginx.conf` to the same directory of `nginx.exe`
1. run `./monitor.exe`