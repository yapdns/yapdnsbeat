# YAPDNSBeat

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
## Running the client

Client by default checks for `yapdnsbeat.yml` config file in current directory

Custom config file

	./yapdnsbeat -c custom.yml

More Config options - common to libbeats

	./yapdnsbeat -h

## Config options

	Detailed in yapdnsbeat.yml


