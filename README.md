# SecuritySpy Command Line Interface (secspy)

Works with SecuritySpy v4 and v5.

## Overview

This app was written as a way to test and provide example usage of the
[securityspy](https://golift.io/securityspy) library.
This app will be updated slowly as I feel like adding new feature.

If you have feature requests or bug reports, please open an Issue.

## macOS

Install with homebrew.
```shell
brew install golift/mugs/secspy
secspy --help
```

## Linux

You can install on Linux with a DEB or RPM.
Those are on the [Releases](https://github.com/davidnewhall/secspy/releases) page,
and available in a packagecloud repo. Install the repo and secspy like this:

```shell
curl -s https://golift.io/repo.sh | sudo bash -s - secspy
```

<a href="https://packagecloud.io"><img src="https://packagecloud.io/images/packagecloud-badge.png" /></a>

## Docker

Or use a Docker image.
```shell
docker pull golift/secspy:stable
docker run golift/secspy:stable --help
```

## Example

Example output from `--help`
```
Usage: secspy [--user <user>] [--pass <pass>] [--url <url>] [-c <cmd>] [-a <arg>]
  -a, --arg string       if cmd supports an argument, pass it here. ie. -c pic -a Porch:/tmp/filename.jpg
  -c, --command string   Command to run. Currently supports: events/callback, cams, pic, vid, trigger, files, download, ptz, arm
  -p, --pass string      Password to authenticate with
  -U, --url string       SecuritySpy URL (default "http://127.0.0.1:8000")
  -u, --user string      Username to authenticate with
  -s, --verify-ssl       Validate SSL certificate if using https
  -v, --version          Print the version and exit
pflag: help requested
```
