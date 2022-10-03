# dca
[![Go report]( http://goreportcard.com/badge/bwmarrin/dca)](http://goreportcard.com/report/bwmarrin/dca) [![Build Status](https://travis-ci.org/bwmarrin/dca.svg?branch=master)](https://travis-ci.org/bwmarrin/dca) 

dca is a command line tool that provides an example implementation of the DCA 
audio format.

This tool accepts raw PCM from stdin and outputs valid DCA0 data on stdout.

**NOTE:** Currently this tool only supports DCA0.  DCA1 will be added later.

You can also pipe the output of this program to create a .dca file for later use.

* See [Discordgo](https://github.com/bwmarrin/discordgo) for Discord API bindings in Go.
* See the [bwmarrin/dca](https://github.com/bwmarrin/dca) for more information on the DCA audio format.

**For help with this program or general Go discussion, please join the [Discord 
Gophers](https://discord.gg/0f1SbxBZjYq9jLBk) chat server.**

## Features
* Stereo Audio
* 48khz Sampling Rate
* 20ms / 1920 byte audio frame size
* Bit-rates from 8 kb/s to 128 kb/s
* Optimization setting for VoIP, Audio, and Low Delay audio


## Getting Started

This assumes you already have a working Go environment, if not please see
[this page](https://golang.org/doc/install) first.

### Installing dca

#### Linux

From a terminal run the following command to download and compile this dca tool.

```sh
go install github.com/bwmarrin/dca/cmd/dca@latest
```

This will use the Go install tool to download the dca package and the opus library 
dependency then compile the tool and install it in your Go bin folder.


### Using dca with ffmpeg

This dca tool only accepts PCM input and one of the easiest ways to get that
is by using ffmpeg tool to convert (nearly) any audio file you have into PCM
and then pipe that into this tool.

Below is an example of using ffmpeg to read a file `test.mp3` then convert that 
into PCM and pipe it into this dca tool and save the result as `test.dca`

This uses the default dca settings, you can of course add arguments to the dca 
command to modify it's default behaviour.

```sh
ffmpeg -i test.mp3 -f s16le -ar 48000 -ac 2 pipe:1 | dca > test.dca
```


## Structure

Here is the structure of a DCA0 frame:

```
| 0 | 1 | 2 - Frame Size |
|---|---|----------------|
| Frame |  Opus encoded  |
| Size  |      data      |
| int16 |                |
```
