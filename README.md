# Chromecast

Implements a few Google Chromecast commands. Other than the basic commands, it also allows you to play media files from your computer.

## Playable Media Content

Can load a local media file or a file hosted on the internet on your chromecast with the following format:

```
Supported Media formats:
    - MP3
    - AVI
    - MKV
    - MP4
    - WebM
    - FLAC
    - WAV
```

If an unknown video file is found, it will use `ffmpeg` to transcode it to MP4 and stream it to the chromecast.

## Play Local Media Files

We are able to play local media files by creating a http server that will stream the media file to the cast device.

## Cast DNS Lookup

A DNS multicast is used to determine the Chromecast and Google Home devices.

The cast DNS entry is also cached, this means that if you pass through the device name, `-n <name>`, or the
device uuid, `-u <uuid>`, the results will be cached and it will connect to the chromecast device instantly.

## Installing

```
$ go get github.com/newbits/go-chromecast
```
