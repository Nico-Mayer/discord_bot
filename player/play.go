package player

import (
	"io"
	"log"
	"time"

	"github.com/ClintonCollins/dca"
	"github.com/nico-mayer/go_discordbot/utils"
)

func (p *Player) Play() {

	if p.PlayerStatus != Resting {
		return
	}

	song := p.dequeue()

	log.Printf("Playing song: %v \n", song.Name)

	encodingSession, err := dca.EncodeFile(song.downloadUrl, p.options)
	if err != nil {
		log.Println("Error encoding from yt url")
		log.Println(err)
		return
	}

	defer encodingSession.Cleanup()

	time.Sleep(250 * time.Millisecond)
	err = p.voiceConn.Speaking(true)
	utils.Check(err)

	done := make(chan error)
	stream := dca.NewStream(encodingSession, p.voiceConn, done)

	p.currentStream = stream
	log.Println("Created stream, waiting on finish or err")
	p.PlayerStatus = Playing

	select {
	case err := <-done:
		// Case 1: Song has finished playing.
		log.Println("Song done")

		// Check if there is an error and it's not the end of file (EOF).
		if err != nil && err != io.EOF {
			p.PlayerStatus = Err
			log.Println(err.Error())
			return
		}

		// Set speaking status to false and exit the select block.
		p.voiceConn.Speaking(false)
		break

	case <-p.SkipInterrupt:
		// Case 2: Song is interrupted.
		log.Println("Song interrupted, stop playing")

		// Set speaking status to false and return from the function.
		p.voiceConn.Speaking(false)
		return
	}

	// Additional code executed regardless of the selected case.
	p.voiceConn.Speaking(false)

	// Check if the queue is empty.
	if len(p.queue) == 0 {
		// If the queue is empty, wait for a short duration and then stop playing.
		time.Sleep(250 * time.Millisecond)
		log.Println("Audio done")
		p.Stop()
		p.PlayerStatus = Resting
		return
	}

	// If the queue is not empty, wait for a short duration and log the next song in the queue.
	time.Sleep(250 * time.Millisecond)
	log.Println("Play next in queue")
	p.PlayerStatus = Resting

	// Start playing the next song in the queue asynchronously.
	go p.Play()
}
