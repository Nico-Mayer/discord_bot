package player

import (
	"log"
)

func (p *Player) Enqueue(song *Song) {
	log.Printf("Queueing song %v", song.Name)
	songString := song.Name
	p.QueueList = append(p.QueueList, songString)
	p.queue <- song
}

func (p *Player) dequeue() *Song {
	p.QueueList = p.QueueList[1:]
	return <-p.queue
}
