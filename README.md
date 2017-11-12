# DCA
[![Go report]( http://goreportcard.com/badge/bwmarrin/dca)](http://goreportcard.com/report/bwmarrin/dca) [![Build Status](https://travis-ci.org/bwmarrin/dca.svg?branch=master)](https://travis-ci.org/bwmarrin/dca) 

`dca` is an audio file format that uses opus audio packets and json metadata.


`dca` files are designed to be easily sent directly to Discord with minimal 
additional processing. `dca` files may also be suitable for any other 
service that accepts Opus audio. 


**For help with this program or general Go discussion, please join the [Discord 
Gophers](https://discord.gg/0f1SbxBZjYq9jLBk) chat server.**


### What's here?

This repository hosts the official specification for `dca` and an example 
implementation of the DCA1 specification.

### Official Specifications
* [DCA0 specification](https://github.com/bwmarrin/dca/wiki/DCA0-specification)
* [DCA1 specification](https://github.com/bwmarrin/dca/wiki/DCA1-specification)


### Implementations of DCA

Each of these implementations have their own unique differences.  It is 
recommended to review and evaluate each of them to decide which one fits your
needs best.
 
| Name                                                       | Lang |
| ---------------------------------------------------------- | ---- |
| [dCa](https://github.com/uppfinnarn/dca)                   | C    |
| [dcad](https://github.com/b1naryth1ef/dcad)                | D    |
| [dca](https://github.com/jonas747/dca)                     | Go   |
| [dca](https://github.com/bwmarrin/dca/tree/master/cmd/dca) | Go   |
