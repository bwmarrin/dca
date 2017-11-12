package main

// This file parses the command line options and fires off the appropriate
// functions to create your dca file.

import (
	"flag"
	"fmt"
	"os"

	"layeh.com/gopus"
)

// Parse command line arguments and setup a couple of variables.
func init() {

	// Opus Encoding Options
	flag.IntVar(&AudioChannels, "ac", 2, "audio channels")
	flag.IntVar(&AudioFrameRate, "ar", 48000, "audio sampling rate")
	flag.IntVar(&AudioFrameSize, "as", 960, "audio frame size can be 960 (20ms), 1920 (40ms), or 2880 (60ms)")
	flag.IntVar(&AudioBitrate, "ab", 64, "audio encoding bitrate in kb/s can be 8 - 128")
	flag.StringVar(&AudioApplication, "aa", "audio", "audio application can be voip, audio, or lowdelay")

	flag.Parse()

	MaxBytes = (AudioFrameSize * AudioChannels) * 2 // max size of opus data
}

func main() {

	//////////////////////////////////////////////////////////////////////////
	// Basic validation
	//////////////////////////////////////////////////////////////////////////

	// Make sure the stdin pipe is open
	fi, err := os.Stdin.Stat()
	if err != nil {
		fmt.Println(err)
		return
	}

	if (fi.Mode() & os.ModeCharDevice) == 0 {
	} else {
		fmt.Println("error: stdin is not a pipe.")
		flag.Usage()
		return
	}

	//////////////////////////////////////////////////////////////////////////
	// Create channels, buffers, and encoder for use
	//////////////////////////////////////////////////////////////////////////

	// Create an Open Encoder to use
	OpusEncoder, err = gopus.NewEncoder(AudioFrameRate, AudioChannels, gopus.Audio)
	if err != nil {
		fmt.Println("NewEncoder Error:", err)
		return
	}

	// Set Opus Encoder Bitrate
	if AudioBitrate < 1 || AudioBitrate > 512 {
		AudioBitrate = 64 // Set to Discord default
	}
	OpusEncoder.SetBitrate(AudioBitrate * 1000)

	// Set Opus Encoder Application
	switch AudioApplication {
	case "voip":
		OpusEncoder.SetApplication(gopus.Voip)
	case "audio":
		OpusEncoder.SetApplication(gopus.Audio)
	case "lowdelay":
		OpusEncoder.SetApplication(gopus.RestrictedLowDelay)
	default:
		OpusEncoder.SetApplication(gopus.Audio)
	}

	// Create channels used by the reader/encoder/writer go routines
	EncodeChan = make(chan []int16, 10)
	OutputChan = make(chan []byte, 10)

	//////////////////////////////////////////////////////////////////////////
	// Start reader, encoder, and writer workers.  These add the DCA0 format
	// audio content to the file.
	//////////////////////////////////////////////////////////////////////////

	WaitGroup.Add(1)
	go reader()

	WaitGroup.Add(1)
	go encoder()

	WaitGroup.Add(1)
	go writer()

	// wait for above goroutines to finish, then exit.
	WaitGroup.Wait()
}
