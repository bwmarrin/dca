package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"
	"runtime"
	"sync"

	"github.com/bwmarrin/discordgo"
)

var (
	run *exec.Cmd
)

func main() {

	// NOTE: All of the below fields are required for this example to work correctly.
	var (
		Email     = flag.String("e", "", "Discord account email.")
		Password  = flag.String("p", "", "Discord account password.")
		GuildID   = flag.String("g", "", "Guild ID")
		ChannelID = flag.String("c", "", "Channel ID")
		Folder    = flag.String("f", "", "Folder of files to play.")
		err       error
	)
	flag.Parse()

	// Connect to Discord
	discord, err := discordgo.New(*Email, *Password)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Open Websocket
	err = discord.Open()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Connect to voice channel.
	// NOTE: Setting mute to false, deaf to true.
	err = discord.ChannelVoiceJoin(*GuildID, *ChannelID, false, true)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Hacky loop to prevent sending on a nil channel.
	// TODO: Find a better way.
	for discord.Voice.Ready == false {
		runtime.Gosched()
	}

	// Start loop and attempt to play all files in the given folder
	fmt.Println("Reading Folder: ", *Folder)
	files, _ := ioutil.ReadDir(*Folder)
	for _, f := range files {
		fmt.Println("PlayAudioFile:", f.Name())
		discord.UpdateStatus(0, f.Name())
		PlayAudioFile(discord.Voice, fmt.Sprintf("%s/%s", *Folder, f.Name()))
	}

	// Close connections
	discord.Voice.Close()
	discord.Close()

	return
}

// PlayAudioFile will play the given filename to the already connected
// Discord voice server/channel.  voice websocket and udp socket
// must already be setup before this will work.
func PlayAudioFile(v *discordgo.Voice, filename string) {

	// Create a shell command "object" to run.
	run = exec.Command("ff2opus", filename)
	stdout, err := run.StdoutPipe()
	if err != nil {
		fmt.Println("StdoutPipe Error:", err)
		return
	}

	// Starts the ff2opus command
	err = run.Start()
	if err != nil {
		fmt.Println("RunStart Error:", err)
		return
	}

	// header "buffer"
	var opuslen uint16

	// Send "speaking" packet over the voice websocket
	v.Speaking(true)

	// Send not "speaking" packet over the websocket when we finish
	defer v.Speaking(false)

	for {

		// read "header" from ff2opus
		err = binary.Read(stdout, binary.LittleEndian, &opuslen)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		if err != nil {
			fmt.Println("error reading from ff2opus stdout :", err)
			return
		}

		// read opus data from ff2opus
		opus := make([]byte, opuslen)
		err = binary.Read(stdout, binary.LittleEndian, &opus)
		if err == io.EOF || err == io.ErrUnexpectedEOF {
			return
		}
		if err != nil {
			fmt.Println("error reading from ff2opus stdout :", err)
			return
		}

		// Send received PCM to the sendPCM channel
		v.OpusSend <- opus
	}
}
