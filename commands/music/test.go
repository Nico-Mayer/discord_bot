package music

/* package music

import (
	"context"
	"encoding/binary"
	"io"
	"os"
	"time"

	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/events"
	"github.com/disgoorg/disgo/voice"
	"github.com/disgoorg/snowflake/v2"
	"github.com/nico-mayer/discordbot/config"
)

var PlayCommand = discord.SlashCommandCreate{
	Name:        "play",
	Description: "play music",
}

func PlayCommandHandler(event *events.ApplicationCommandInteractionCreate) {
	client := event.Client()
	conn := client.VoiceManager().CreateConn(config.GUILD_ID)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	if err := conn.Open(ctx, snowflake.MustParse("1082979754312994880"), false, false); err != nil {
		panic("error connecting to voice channel: " + err.Error())
	}
	defer func() {
		closeCtx, closeCancel := context.WithTimeout(context.Background(), time.Second*10)
		defer closeCancel()
		conn.Close(closeCtx)
	}()

	if err := conn.SetSpeaking(ctx, voice.SpeakingFlagMicrophone); err != nil {
		panic("error setting speaking flag: " + err.Error())
	}
	writeOpus(conn.UDP())
}

func writeOpus(w io.Writer) {
	file, err := os.Open("nico.dca")
	if err != nil {
		panic("error opening file: " + err.Error())
	}
	ticker := time.NewTicker(time.Millisecond * 20)
	defer ticker.Stop()

	var lenBuf [4]byte
	for range ticker.C {
		_, err = io.ReadFull(file, lenBuf[:])
		if err != nil {
			if err == io.EOF {
				_ = file.Close()
				return
			}
			panic("error reading file: " + err.Error())
		}

		// Read the integer
		frameLen := int64(binary.LittleEndian.Uint32(lenBuf[:]))

		// Copy the frame.
		_, err = io.CopyN(w, file, frameLen)
		if err != nil && err != io.EOF {
			_ = file.Close()
			return
		}
	}
}
*/
