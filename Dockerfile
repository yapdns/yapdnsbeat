FROM golang:1.6

ENV GLIDE_VERSION 0.10.2

RUN apt-get update \
 	&& apt-get install -y unzip --no-install-recommends \
	&& rm -rf /var/lib/apt/lists/*

ENV GLIDE_DOWNLOAD_URL https://github.com/Masterminds/glide/releases/download/$GLIDE_VERSION/glide-$GLIDE_VERSION-linux-amd64.zip

RUN curl -fsSL "$GLIDE_DOWNLOAD_URL" -o glide.zip \
	&& unzip glide.zip  linux-amd64/glide \
	&& mv linux-amd64/glide /usr/local/bin \
	&& rm -rf linux-amd64 \
	&& rm glide.zip

RUN mkdir -p /go/src/github.com/yapdns/yapdns-client
ADD . /go/src/github.com/yapdns/yapdns-client
WORKDIR /go/src/github.com/yapdns/yapdns-client
RUN glide install
RUN go build && ./yapdns-client -e -v
