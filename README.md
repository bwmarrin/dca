dca  [![Go report](http://goreportcard.com/badge/bwmarrin/dca)](http://goreportcard.com/report/bwmarrin/dca) [![Build Status](https://travis-ci.org/bwmarrin/discordgo.svg?branch=master)](https://travis-ci.org/bwmarrin/dca)
====

dca is a command line tool that wraps ffmpeg to create opus audio data suitable
for use with the [Discord](https://discordapp.com/) chat software.

If you are developing a library for use with Discord you can use this program
as a way to generate the opus audio data from any standard audio file.

You can also pipe the output of this program to create a .dca file for later use.

* See [Discordgo](https://github.com/bwmarrin/discordgo) for Discord API bindings in Go.

Join [#go_discordgo](https://discord.gg/0SBTUU1wZTWT6sqd) Discord chat channel 
for support.

## Features
* Sampling rates from 8 to 48khz
* Bit-rates from 6 kb/s to 510 kb/s
* Support CBR and VBR support.
* Support mono and stereo
* Frame sizes from 10ms to 60ms


## Getting Started

### Installing

dca has been tested to compile on FreeBSD 10 (Go 1.5.1), OS X 10.10, Windows 10.


### Windows
Provided by Axiom :) -- Very ROUGH DRAFT
```
Install Go for Windows
Setup gopath to some empty folder (for example, I made mine C:\gopath)
Install winbuilds (http://win-builds.org/doku.php) (handles our external dependencies for us including 64bit gcc needed to compile)
Inside of winbuilds, install gcc, its dependencies, and opus (might be under libopus). If you're really unsure, just hit process in the top right which will install everything.
Open cmd and cd into dca repository directory
Run go build
???
Profit!
```

### OS X
Provided by Uniquoooo :) -- Very ROUGH DRAFT.
```
1. get homebrew
2. brew install ffmpeg
3. brew install opus
4. brew install golang
5. go get github.com/bwmarrin/dca
```


### Usage

```
Usage of ./dca:
  -ac int
    audio channels (default 2)
  -ar int
    audio sampling rate (default 48000)
  -i string
    infile (default pipe:0)
```

You may also pass pipe pcm16 audio into dca instead of providing an input file.


## Examples

See the example folder.


## Contributing

While contributions are always welcome - this code is in a very early and 
incomplete stage and massive changes, including entire re-writes, could still
happen.  In other words, probably not worth your time right now :)

## List of Discord APIs

See [this chart](https://abal.moe/Discord/Libraries.html) for a feature 
comparison and list of other Discord API libraries.
