# dca
[![Go report]( http://goreportcard.com/badge/bwmarrin/dca)](http://goreportcard.com/report/bwmarrin/dca) [![Build Status](https://travis-ci.org/bwmarrin/dca.svg?branch=master)](https://travis-ci.org/bwmarrin/dca) 

dca is a command line tool that provides an example implementation of the DCA 
audio format.

This tool accepts json metadata from a file, and raw PCM from stdin and outputs 
valid DCA0 data on stdout.

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


## Structure

Here is the structure of a DCA1 file header:

```
| 0 | 1 | 2 |         3        |  4  |  5  |  6  |  7  | 8 - JSON Size |
|---|---|---|------------------|-----------------------|---------------|
|    DCA    |  Version Number  |       JSON Size       | JSON Metadata |
|  Magic Header with Version   |      signed int32     |               |
```

Here is the structure of A DCA0/1 frame:

```
| 0 | 1 | 2 - Frame Size |
|---|---|----------------|
| Frame |  Opus encoded  |
| Size  |      data      |
| int16 |                |
```
