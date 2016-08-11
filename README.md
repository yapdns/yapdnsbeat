# YAPDNS Client

Watch log files, extract DNS records and send them to YAPDNS application

## Installation

Install and configure [Go](https://golang.org/doc/install).

Install and update this go package with:

```bash
go get -u github.com/yapdns/yapdnsbeat
```

To create `yapdnsbeat` binary
```bash
cd $GOPATH/github.com/yapdns/yapdnsbeat
# build dev branch
git checkout dev
go build
```
