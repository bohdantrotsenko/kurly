# kurly

[![Build Status](https://gitlab.com/davidjpeacock/kurly/badges/master/build.svg)](https://gitlab.com/davidjpeacock/kurly/-/jobs)

kurly is an alternative to the widely popular curl program.

kurly is designed to operate in a similar manner to curl, with select features.
Notably, kurly is not aiming for feature parity, but common flags and mechanisms
particularly within the HTTP(S) realm are to be expected.

The current authors are not security experts, but want to contribute to the fledging
movement of replacing key tools and services with equivalents based on modern
and safe languages.  We recognize that people are fallible (including
ourselves), and for this reason believe we need all the help we can get.

Several languages exist which could be used to fulfill our goal, but in this case
we picked Golang.

## Installation

**Pre-requisite: Golang version 1.7.4 or higher.**

From source you can simply:

`go get gitlab.com/davidjpeacock/kurly`

## OS Distribution Package

`kurly` can be installed through package management systems on the following platforms:

* Arch Linux via Arch User Repos
  + For stable versions : `pacaur -S kurly` or `yaourt -S kurly`
  + For tip/development versions : `pacaur -S kurly-git` or `yaourt -S kurly-git` 

*If you're a package maintainer and you have prepared kurly for your OS of choice, please
PR this section.*

## Binary download

Binaries are provided for the following platforms:

* [Mac OS X amd64](https://gitlab.com/davidjpeacock/kurly/-/jobs/artifacts/master/raw/artifacts/darwin-amd64/kurly-darwin-amd64.tar.gz?job=compile)
* [Windows amd64](https://gitlab.com/davidjpeacock/kurly/-/jobs/artifacts/master/raw/artifacts/windows-amd64/kurly-windows-amd64.zip?job=compile)
* [Linux amd64](https://gitlab.com/davidjpeacock/kurly/-/jobs/artifacts/master/raw/artifacts/linux-amd64/kurly-linux-amd64.tar.gz?job=compile)
* [Linux arm](https://gitlab.com/davidjpeacock/kurly/-/jobs/artifacts/master/raw/artifacts/linux-arm/kurly-linux-arm.tar.gz?job=compile)
* [FreeBSD amd64](https://gitlab.com/davidjpeacock/kurly/-/jobs/artifacts/master/raw/artifacts/freebsd-amd64/kurly-freebsd-amd64.tar.gz?job=compile)
* [FreeBSD arm](https://gitlab.com/davidjpeacock/kurly/-/jobs/artifacts/master/raw/artifacts/freebsd-arm/kurly-freebsd-arm.tar.gz?job=compile)
* [OpenBSD amd64](https://gitlab.com/davidjpeacock/kurly/-/jobs/artifacts/master/raw/artifacts/openbsd-amd64/kurly-openbsd-amd64.tar.gz?job=compile)
* [OpenBSD arm](https://gitlab.com/davidjpeacock/kurly/-/jobs/artifacts/master/raw/artifacts/openbsd-arm/kurly-openbsd-arm.tar.gz?job=compile)
* [NetBSD amd64](https://gitlab.com/davidjpeacock/kurly/-/jobs/artifacts/master/raw/artifacts/netbsd-amd64/kurly-netbsd-amd64.tar.gz?job=compile)
* [NetBSD arm](https://gitlab.com/davidjpeacock/kurly/-/jobs/artifacts/master/raw/artifacts/netbsd-arm/kurly-netbsd-arm.tar.gz?job=compile)

## Usage

See `kurly --help` for usage information.

## Examples

Verbose output, showing headers
```
$ kurly -v https://httpbin.org/ip
*   Trying 23.23.171.5...
* TCP_NODELAY set
* Connected to httpbin.org (23.23.171.5) port 443 (#0)
* APLN, server accepted to use http/1.1
* TLSv1.2, TLS Handshake finished
* SSL connection using TLSv1.2 / ECDHE-RSA-AES-128-GCM-SHA256
* Server certificate:
*  subject: CN=httpbin.org
*  start date: Thu, 11 Jan 2018 23:37:29 UTC
*  expire date: Wed, 11 Apr 2018 23:37:29 UTC
*  issuer: C=US; O=Let's Encrypt; CN=Let's Encrypt Authority X3
*  SSL certificate verify ok.
> GET /ip HTTP/1.1
> User-Agent [Kurly/master]
> Accept [*/*]
> Host [httpbin.org]
< HTTP/1.1 200 OK
< Server [meinheld/0.6.1]
< Content-Type [application/json]
< Access-Control-Allow-Credentials [true]
< Content-Length [33]
< Via [1.1 vegur]
< Connection [keep-alive]
< Date [Mon, 12 Mar 2018 19:18:11 GMT]
< Access-Control-Allow-Origin [*]
< X-Powered-By [Flask]
< X-Processed-Time [0]
[<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<<] 33 B/33 B
{
  "origin": "43.122.23.223"
}
```

Download file, preserving remote filename, timestamp, and following redirects
```
$ kurly -R -O -L http://cdimage.debian.org/debian-cd/current/amd64/iso-cd/debian-8.7.1-amd64-netinst.iso
[<<<<<<                                ] 41.2 MB/260 MB
```

Upload file
```
$ kurly -T ~/Downloads/image.jpeg https://httpbin.org/put
```

Posting elements with -d
```
$ kurly -d bingo=bongo https://httpbin.org/post
```

## Roadmap

Succinctly, we're planning to cover all curl-like features relevant to HTTP(S), and would
love you to help.

## Contributing

Bug reports, feature requests, and pull requests are all welcome.  Thank you!

Please see [CONTRIBUTING.md](https://gitlab.com/davidjpeacock/kurly/blob/master/CONTRIBUTING.md) for details of how to work with us.

## Maintainers

kurly is brought to you and maintained by:

* [Akilan Elango](https://gitlab.com/aki237)
* [Al S-M](https://github.com/alsm)
* [David J Peacock](https://gitlab.com/davidjpeacock)

## License

kurly is Copyright (c) 2017-2018 David J Peacock and Al S-M, and is published under the Apache 2.0 license.
