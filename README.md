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

## JSON Structure

Here is the structure of the JSON metadata:

```
{
    "dca": { // Contains information about the DCA encoder
        "version": 1, // The DCA format version that the file is encoded with.
        "tool": { // Information about the tool that encoded the DCA file.
            "name": "dca-encoder", // The name of the tool.
            "version": "1.0.0", // The version of the tool.
            "rev": "bwmarrin/dca#32361ee92fcbd0e404b2be18adf497a45fef4a5f", // The Git revision of the tool.
            "url": "https://github.com/bwmarrin/dca/", // A URL to the tool.
            "author": "bwmarrin" // The author of the tool.
        }
    },
    "info": { // Information about the song from FFmpeg. Most of this is obvious.
        "title": "Out of Control", 
        "artist": "Nothing's Carved in Stone",
        "album": "Revolt",
        "genre": "jrock",
        "comments": "Second Opening for the anime Psycho Pass",
        "cover": "" // A Base64 encoded JPEG image of the songs cover art.
    },
    "origin": { // Information about where the song came from.
        "source": "file", // Whether it was streamed, from a file etc.
        "abr": 192000, // The original bitrate of the file.
        "channels": 2, // The original amount of channels of the file.
        "encoding": "MP3/MPEG-2L3", // The original encoding of the file.
        "url": "https://www.dropbox.com/s/bwc73zb44o3tj3m/Out%20of%20Control.mp3?dl=0" // A URL or path to the file.
    },
    "opus": { // Information about the Opus encoder.
        "abr": 64000, // The bitrate the opus was encoded with.
        "sample_rate": 48000, // The sample rate the opus was encoded with.
        "mode": "voip", // The application mode the opus was encoded with.
        "frame_size": 960, // The frame size the opus was encoded with.
        "channels": 2 // The amount of channels the opus was encoded with.
    }
}
```