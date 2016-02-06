package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"

	"github.com/layeh/gopus"
)

var (
	channels    int = 2                   // 1 for mono, 2 for stereo
	frameRate   int = 48000               // audio sampling rate
	frameSize   int = 960                 // uint16 size of each audio frame
	maxBytes    int = (frameSize * 2) * 2 // max size of opus data
	opusEncoder *gopus.Encoder
	run         *exec.Cmd
)

// very simple program that wraps ffmpeg and outputs opus data
func main() {

	var err error

	if len(os.Args) < 2 {
		fmt.Println("Must supply the filename to process.")
		return
	}

	filename := os.Args[1]

	opusEncoder, err = gopus.NewEncoder(frameRate, channels, gopus.Audio)
	if err != nil {
		fmt.Println("NewEncoder Error:", err)
		return
	}

	// Create a shell command "object" to run.
	run = exec.Command("ffmpeg", "-i", filename, "-f", "s16le", "-ar", strconv.Itoa(frameRate), "-ac", strconv.Itoa(channels), "pipe:1")
	stdout, err := run.StdoutPipe()
	if err != nil {
		fmt.Println("StdoutPipe Error:", err)
		return
	}

	// Starts the ffmpeg command
	err = run.Start()
	if err != nil {
		fmt.Println("RunStart Error:", err)
		return
	}

	// buffer used during loop below
	audiobuf := make([]int16, frameSize*channels)

	// "header" :)
	var opuslen uint16

	for {

		// read data from ffmpeg stdout
		err = binary.Read(stdout, binary.LittleEndian, &audiobuf)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		if err != nil {
			fmt.Println("error reading from ffmpeg stdout :", err)
			return
		}

		// try encoding pcm frame with Opus
		opus, err := opusEncoder.Encode(audiobuf, frameSize, maxBytes)
		if err != nil {
			fmt.Println("Encoding Error:", err)
			return
		}

		// write header
		opuslen = uint16(len(opus))
		err = binary.Write(os.Stdout, binary.LittleEndian, &opuslen)
		if err != nil {
			fmt.Println("error writing output: ", err)
			return
		}

		// write opus data to stdout
		err = binary.Write(os.Stdout, binary.LittleEndian, &opus)
		if err != nil {
			fmt.Println("error writing output: ", err)
			return
		}
	}
}
