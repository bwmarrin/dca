package main

import (
	"bytes"
	"bufio"
	"encoding/binary"
	"encoding/json"
	"encoding/base64"
	"flag"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strconv"
	"sync"
	"image/png"
	"image/jpeg"

	"github.com/layeh/gopus"
)

// Define constants
const (
	// The current version of the DCA format
	FormatVersion int8 = 1

	// The current version of the DCA program
	ProgramVersion string = "0.0.1"

	// The URL to the GitHub repository of DCA
	GitHubRepositoryURL string = "https://github.com/bwmarrin/dca"
)

// All global variables used within the program
var (
	// Buffer for some commands
	CmdBuf bytes.Buffer
	PngBuf bytes.Buffer

	// Metadata structures
	Metadata	MetadataStruct
	FFprobeData FFprobeMetadata

	// Magic bytes to write at the start of a DCA file
	MagicBytes string = fmt.Sprintf("DCA%d", FormatVersion)

	// 1 for mono, 2 for stereo
	Channels int

	// Must be one of 8000, 12000, 16000, 24000, or 48000.
	// Discord only uses 48000 currently.
	FrameRate int

	// Rates from 500 to 512000 bits per second are meaningful
	// Discord only uses 8000 to 128000 and default is 64000
	Bitrate int

	// Must be one of voip, audio, or lowdelay.
	// DCA defaults to audio which is ideal for music
	// Not sure what Discord uses here, probably voip
	Application string

	FrameSize int // uint16 size of each audio frame
	MaxBytes  int // max size of opus data

	Volume int // change audio volume (256=normal)

	OpusEncoder *gopus.Encoder

	InFile string
	CoverFormat string = "jpeg"

	OutFile string = "pipe:1"
	OutBuf  []byte

	EncodeChan chan []int16
	OutputChan chan []byte

	err error

	wg sync.WaitGroup
)

// init configures and parses the command line arguments
func init() {

	flag.StringVar(&InFile, "i", "pipe:0", "infile")
	flag.IntVar(&Volume, "vol", 256, "change audio volume (256=normal)")
	flag.IntVar(&Channels, "ac", 2, "audio channels")
	flag.IntVar(&FrameRate, "ar", 48000, "audio sampling rate")
	flag.IntVar(&FrameSize, "as", 960, "audio frame size can be 960 (20ms), 1920 (40ms), or 2880 (60ms)")
	flag.IntVar(&Bitrate, "ab", 64, "audio encoding bitrate in kb/s can be 8 - 128")
	flag.StringVar(&Application, "aa", "audio", "audio application can be voip, audio, or lowdelay")
	flag.StringVar(&CoverFormat, "format", "jpeg", "format the cover art will be encoded with")

	if len(os.Args) < 2 {
		flag.Usage()
		os.Exit(1)
	}

	flag.Parse()

	MaxBytes = (FrameSize * Channels) * 2 // max size of opus data
}

// very simple program that wraps ffmpeg and outputs raw opus data frames
// with a uint16 header for each frame with the frame length in bytes
func main() {

	//////////////////////////////////////////////////////////////////////////
	// BLOCK : Basic setup and validation
	//////////////////////////////////////////////////////////////////////////

	// If only one argument provided assume it's a filename.
	if len(os.Args) == 2 {
		InFile = os.Args[1]
	}

	// If reading from a file, verify it exists.
	if InFile != "pipe:0" {

		if _, err := os.Stat(InFile); os.IsNotExist(err) {
			fmt.Println("error: infile does not exist")
			flag.Usage()
			return
		}
	}

	// If reading from pipe, make sure pipe is open
	if InFile == "pipe:0" {
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
	}

	//////////////////////////////////////////////////////////////////////////
	// BLOCK : Create chans, buffers, and encoder for use
	//////////////////////////////////////////////////////////////////////////

	// create an opusEncoder to use
	OpusEncoder, err = gopus.NewEncoder(FrameRate, Channels, gopus.Audio)
	if err != nil {
		fmt.Println("NewEncoder Error:", err)
		return
	}

	// set opus encoding options
	//	OpusEncoder.SetVbr(true)                // bool

	if Bitrate < 1 || Bitrate > 512 {
		Bitrate = 64 // Set to Discord default
	}
	OpusEncoder.SetBitrate(Bitrate * 1000)

	switch Application {
	case "voip":
		OpusEncoder.SetApplication(gopus.Voip)
	case "audio":
		OpusEncoder.SetApplication(gopus.Audio)
	case "lowdelay":
		OpusEncoder.SetApplication(gopus.RestrictedLowDelay)
	default:
		OpusEncoder.SetApplication(gopus.Audio)
	}

	OutputChan = make(chan []byte, 10)
	EncodeChan = make(chan []int16, 10)

	// Setup the metadata
	Metadata = MetadataStruct{
		Dca: &DCAMetadata{
			Version: FormatVersion,
			Tool: &DCAToolMetadata{
				Name: "dca",
				Version: ProgramVersion,
				Revision: "",
				Url: GitHubRepositoryURL,
				Author: "bwmarrin",
			},
		},
		SongInfo: &SongMetadata{},
		Origin: &OriginMetadata{},
		Opus: &OpusMetadata{
			Bitrate: Bitrate * 1000,
			SampleRate: FrameRate,
			Application: Application,
			FrameSize: FrameSize,
			Channels: Channels,
		},
	}
	_ = Metadata

	// try get the git revision
	git := exec.Command("cd $GOPATH/src/github.com/bwmarrin/dca && git rev-parse HEAD")
	git.Stdout = &CmdBuf

	err = git.Start()
	if err == nil {
		err = git.Wait()
		if err != nil {
			fmt.Println("Git Error:", err)
			return
		}

		Metadata.Dca.Tool.Revision = CmdBuf.String()
	}

	CmdBuf.Reset()

	// get ffprobe data
	if InFile != "pipe:0" {
		ffprobe := exec.Command("ffprobe", "-v", "quiet", "-print_format", "json", "-show_format", InFile)
		ffprobe.Stdout = &CmdBuf

		err = ffprobe.Start()
		if err != nil {
			fmt.Println("RunStart Error:", err)
			return
		}

		err = ffprobe.Wait()
		if err != nil {
			fmt.Println("FFprobe Error:", err)
			return
		}

		err = json.Unmarshal(CmdBuf.Bytes(), &FFprobeData)
		if err != nil {
			fmt.Println("Erorr unmarshaling the FFprobe JSON:", err)
			return
		}

		Metadata.SongInfo = &SongMetadata{
			Title: FFprobeData.Format.Tags.Title,
			Artist: FFprobeData.Format.Tags.Artist,
			Album: FFprobeData.Format.Tags.Album,
			Genre: FFprobeData.Format.Tags.Genre,
			Comments: "", // change later?
		}

		Metadata.Origin = &OriginMetadata{
			Source: "file",
			Bitrate: FFprobeData.Format.Bitrate,
			Channels: Channels,
			Encoding: FFprobeData.Format.FormatLongName,
			Url: FFprobeData.Format.FileName,
		}

		CmdBuf.Reset()

		// get cover art
		cover := exec.Command("ffmpeg", "-loglevel", "0", "-i", InFile, "-f", "singlejpeg", "pipe:1")
		cover.Stdout = &CmdBuf

		err = cover.Start()
		if err != nil {
			fmt.Println("RunStart Error:", err)
			return
		}

		err = cover.Wait()
		if err == nil {
			buf := bytes.NewBufferString(CmdBuf.String())

			if CoverFormat == "png" {
				img, err := jpeg.Decode(buf)
				if err == nil { // silently drop it, no image
					err = png.Encode(&PngBuf, img)
					if err == nil {
						Metadata.SongInfo.Cover = base64.StdEncoding.EncodeToString(PngBuf.Bytes())
					}
				}
			} else {
				encodedImage := base64.StdEncoding.EncodeToString(CmdBuf.Bytes())
				Metadata.SongInfo.Cover = encodedImage
			}
		}

		CmdBuf.Reset()
		PngBuf.Reset()
	}

	//////////////////////////////////////////////////////////////////////////
	// BLOCK : Start reader and writer workers
	//////////////////////////////////////////////////////////////////////////

	wg.Add(1)
	go reader()

	wg.Add(1)
	go encoder()

	wg.Add(1)
	go writer()

	// wait for above goroutines to finish, then exit.
	wg.Wait()
}

// reader reads from the input
func reader() {

	defer func() {
		close(EncodeChan)
		wg.Done()
	}()

	// read from file
	if InFile != "pipe:0" {

		// Create a shell command "object" to run.
		ffmpeg := exec.Command("ffmpeg", "-i", InFile, "-vol", strconv.Itoa(Volume), "-f", "s16le", "-ar", strconv.Itoa(FrameRate), "-ac", strconv.Itoa(Channels), "pipe:1")
		stdout, err := ffmpeg.StdoutPipe()
		if err != nil {
			fmt.Println("StdoutPipe Error:", err)
			return
		}

		// Starts the ffmpeg command
		err = ffmpeg.Start()
		if err != nil {
			fmt.Println("RunStart Error:", err)
			return
		}

		for {

			// read data from ffmpeg stdout
			InBuf := make([]int16, FrameSize*Channels)
			err = binary.Read(stdout, binary.LittleEndian, &InBuf)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return
			}
			if err != nil {
				fmt.Println("error reading from ffmpeg stdout :", err)
				return
			}

			// write pcm data to the EncodeChan
			EncodeChan <- InBuf

		}
	}

	// read input from stdin pipe
	if InFile == "pipe:0" {

		// 16KB input buffer
		rbuf := bufio.NewReaderSize(os.Stdin, 16384)
		for {

			// read data from stdin
			InBuf := make([]int16, FrameSize*Channels)

			err = binary.Read(rbuf, binary.LittleEndian, &InBuf)
			if err == io.EOF || err == io.ErrUnexpectedEOF {
				return
			}
			if err != nil {
				fmt.Println("error reading from ffmpeg stdout :", err)
				return
			}

			// write pcm data to the EncodeChan
			EncodeChan <- InBuf
		}
	}

}

// encoder listens on the EncodeChan and encodes provided PCM16 data
// to opus, then sends the encoded data to the OutputChan
func encoder() {

	defer func() {
		close(OutputChan)
		wg.Done()
	}()

	for {
		pcm, ok := <-EncodeChan
		if !ok {
			// if chan closed, exit
			return
		}

		// try encoding pcm frame with Opus
		opus, err := OpusEncoder.Encode(pcm, FrameSize, MaxBytes)
		if err != nil {
			fmt.Println("Encoding Error:", err)
			return
		}

		// write opus data to OutputChan
		OutputChan <- opus
	}
}

// writer listens on the OutputChan and writes the output to stdout pipe
// TODO: Add support for writing directly to a file
func writer() {

	defer wg.Done()

	var opuslen int16
	var jsonlen int32

	// 16KB output buffer
	wbuf := bufio.NewWriterSize(os.Stdout, 16384)

	// write the magic bytes
	fmt.Print(MagicBytes)

	// encode and write json length
	json, err := json.Marshal(Metadata)
	if err != nil {
		fmt.Println("Failed to encode the Metadata JSON:", err)
		return
	}

	jsonlen = int32(len(json))
	err = binary.Write(wbuf, binary.LittleEndian, &jsonlen)
	if err != nil {
		fmt.Println("error writing output: ", err)
		return
	}

	// write the actual json
	wbuf.Write(json)

	for {
		opus, ok := <-OutputChan
		if !ok {
			// if chan closed, exit
			return
		}

		// write header
		opuslen = int16(len(opus))
		err = binary.Write(wbuf, binary.LittleEndian, &opuslen)
		if err != nil {
			fmt.Println("error writing output: ", err)
			return
		}

		// write opus data to stdout
		err = binary.Write(wbuf, binary.LittleEndian, &opus)
		if err != nil {
			fmt.Println("error writing output: ", err)
			return
		}
	}
}
