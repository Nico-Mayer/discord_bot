package player

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/ClintonCollins/dca"
	"github.com/bwmarrin/discordgo"
)

type PlayerStatus int32

const (
	Resting PlayerStatus = 0
	Playing PlayerStatus = 1
	Paused  PlayerStatus = 2
	Err     PlayerStatus = 3
)

type Player struct {
	Session       *discordgo.Session
	voiceConn     *discordgo.VoiceConnection
	queue         chan *Song
	queueList     []string
	SkipInterrupt chan bool
	currentStream *dca.StreamingSession
	PlayerStatus  PlayerStatus
}

func NewPlayer(s *discordgo.Session) *Player {
	return &Player{
		Session:       s,
		queue:         make(chan *Song, 100),
		SkipInterrupt: make(chan bool, 1),
		PlayerStatus:  Resting,
	}
}

func (p *Player) Play(query string, voiceState *discordgo.VoiceState) {
	p.dequeue()

	song, err := GetSongInfo(query)
	log.Printf("Playing song: %v \n", song.name)
	if err != nil {
		fmt.Println(err)
		return
	}

	var options = &dca.EncodeOptions{
		Volume:           200,
		Channels:         2,
		FrameRate:        48000,
		FrameDuration:    20,
		Bitrate:          64,
		Application:      dca.AudioApplicationLowDelay,
		CompressionLevel: 10,
		PacketLoss:       3,
		BufferedFrames:   200,
		VBR:              true,
		StartTime:        0,
		VolumeFloat:      0.5,
		RawOutput:        true,
	}

	encodingSession, err := dca.EncodeFile(song.downloadUrl, options)

	if err != nil {
		log.Println("Error encoding from yt url")
		log.Println(err)
		return
	}

	defer encodingSession.Cleanup()

	err = p.joinChannel(voiceState)
	if err != nil {
		log.Println("Error joining voice channel")
		log.Println(err)
		return
	}
	time.Sleep(250 * time.Millisecond)
	err = p.voiceConn.Speaking(true)
	if err != nil {
		fmt.Println(err)
	}

	done := make(chan error)
	stream := dca.NewStream(encodingSession, p.voiceConn, done)

	p.currentStream = stream
	log.Println("Created stream, waiting on finish or err")
	p.PlayerStatus = Playing

	/* log.Printf("EncodingSession:  #%v", encodingSession)
	log.Printf("Stream: #%v", stream) */

	err = <-done
	if err != nil && err != io.EOF {
		fmt.Println(err)
	}

}

func (p *Player) dequeue() *Song {
	if len(p.queueList) == 0 {
		return nil
	}
	p.queueList = p.queueList[1:]
	return <-p.queue
}

func (p *Player) stop() {
	p.voiceConn.Disconnect()
	p.voiceConn = nil
	p.PlayerStatus = Resting
}

func (p *Player) joinChannel(vs *discordgo.VoiceState) error {
	if p.voiceConn != nil {
		p.voiceConn.Disconnect()
	}

	if p.voiceConn == nil {
		conn, err := p.Session.ChannelVoiceJoin(vs.GuildID, vs.ChannelID, false, true)
		if err != nil {
			return err
		}
		p.voiceConn = conn
	}
	return nil
}
