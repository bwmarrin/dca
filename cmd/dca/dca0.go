// This file contains the code for the dca0 spec format
package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"sync"

	"layeh.com/gopus"
)

var (

	// 1 for mono, 2 for stereo
	AudioChannels int

	// Must be one of 8000, 12000, 16000, 24000, or 48000.
	// Discord only uses 48000 currently.
	AudioFrameRate int

	// Rates from 500 to 512000 bits per second are meaningful
	// Discord only uses 8000 to 128000 and default is 64000
	AudioBitrate int

	// Must be one of voip, audio, or lowdelay.
	// DCA defaults to audio which is ideal for music
	// Not sure what Discord uses here, probably voip
	AudioApplication string

	// uint16 size of each audio frame
	AudioFrameSize int

	// max size of opus data
	MaxBytes int

	OpusEncoder *gopus.Encoder
	EncodeChan  chan []int16
	OutputChan  chan []byte
	WaitGroup   sync.WaitGroup
)

// reader reads from the input
func reader() {

	var err error

	defer func() {
		close(EncodeChan)
		WaitGroup.Done()
	}()

	// Create a 16KB input buffer
	stdin := bufio.NewReaderSize(os.Stdin, 16384)

	// Loop over the stdin input and pass the data to the encoder.
	for {

		buf := make([]int16, AudioFrameSize*AudioChannels)

		err = binary.Read(stdin, binary.LittleEndian, &buf)
		if err == io.EOF {
			// Okay! There's nothing left, time to quit.
			return
		}

		if err == io.ErrUnexpectedEOF {
			// Well there's just a tiny bit left, lets encode it, then quit.
			EncodeChan <- buf
			return
		}

		if err != nil {
			// Oh no, something went wrong!
			log.Println("error reading from stdin,", err)
			return
		}

		// write pcm data to the EncodeChan
		EncodeChan <- buf
	}

}

// encoder listens on the EncodeChan and encodes provided PCM16 data
// to opus, then sends the encoded data to the OutputChan
func encoder() {

	defer func() {
		close(OutputChan)
		WaitGroup.Done()
	}()

	for {
		pcm, ok := <-EncodeChan
		if !ok {
			// if chan closed, exit
			return
		}

		// try encoding pcm frame with Opus
		opus, err := OpusEncoder.Encode(pcm, AudioFrameSize, MaxBytes)
		if err != nil {
			fmt.Println("Encoding Error:", err)
			return
		}

		// write opus data to OutputChan
		OutputChan <- opus
	}
}

// writer listens on the OutputChan and writes the output to stdout pipe
func writer() {

	defer WaitGroup.Done()

	var opuslen int16
	var err error

	// 16KB output buffer
	stdout := bufio.NewWriterSize(os.Stdout, 16384)
	defer stdout.Flush()

	for {
		opus, ok := <-OutputChan
		if !ok {
			// if chan closed, exit
			return
		}

		// write header
		opuslen = int16(len(opus))
		err = binary.Write(stdout, binary.LittleEndian, &opuslen)
		if err != nil {
			fmt.Println("error writing output: ", err)
			return
		}

		// write opus data to stdout
		err = binary.Write(stdout, binary.LittleEndian, &opus)
		if err != nil {
			fmt.Println("error writing output: ", err)
			return
		}
	}
}
