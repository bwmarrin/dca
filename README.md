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
* Stereo Audio
* 48khz Sampling Rate
* 20ms / 1920 byte audio frame size
* Bit-rates from 8 kb/s to 128 kb/s
* Optimization setting for VoIP, Audio, and Low Delay audio


## Getting Started

### Installing

dca has been tested to compile on FreeBSD 10 (Go 1.5.1), OS X 10.10, Windows 10.

### Ubuntu 14.04.3 LTS
Provided by Uniquoooo
```
# basics
sudo apt-get update
sudo apt-get install golang git gcc make pkg-config --yes
# golang
mkdir $HOME/go
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc
# ffmpeg
sudo add-apt-repository ppa:kirillshkrogalev/ffmpeg-next
sudo apt-get update
sudo apt-get install ffmpeg --yes
# opus
wget http://downloads.xiph.org/releases/opus/opus-1.1.2.tar.gz
tar -zxvf opus-1.1.2.tar.gz
cd opus-1.1.2
./configure
make && sudo make install
cd ../
rm -r opus-1.1.2 opus-1.1.2.tar.gz
# install dca
go get github.com/bwmarrin/dca
```


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

### Windows (Pacman)
Provided by iopred.
First, install msys2 then install pacman
```
$ pacman -S mingw64/mingw-w64-x86_64-pkg-config
$ pacman -S mingw64/mingw-w64-x86_64-opusfile
$ go get github.com/bwmarrin/dca
$ go install github.com/bwmarrin/dca
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
  -aa string
        audio application can be voip, audio, or lowdelay (default "audio")
  -ab int
        audio encoding bitrate in kb/s can be 8 - 128 (default 64)
  -ac int
        audio channels (default 2)
  -ar int
        audio sampling rate (default 48000)
  -as int
        audio frame size can be 960 (20ms), 1920 (40ms), or 2880 (60ms) (default 960)
  -cf string
        format the cover art will be encoded with (default "jpeg")
  -i string
        infile (default "pipe:0")
  -vol int
        change audio volume (256=normal) (default 256)
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

## File Structure

Here is the structure of a DCA file header:

```
| 0 | 1 | 2 |         3        |  4  |  5  |  6  |  7  | 8 - JSON Size |
|---|---|---|------------------|-----------------------|---------------|
|    DCA    |  Version Number  |       JSON Size       | JSON Metadata |
|  Magic Header with Version   |      signed int32     |               |
```

Here is the structure of A DCA frame:

```
| 0 | 1 | 2 - Frame Size |
|---|---|----------------|
| Frame |  Opus encoded  |
| Size  |      data      |
| int16 |                |
```
