package player

import (
	"fmt"
	"log"
)

func (p *Player) Enqueue(song *Song) {
	log.Printf("Queueing song %v", song.Name)
	songString := fmt.Sprintf("-- :%v \n", song.Name)
	p.queueList = append(p.queueList, songString)
	p.queue <- song
}

func (p *Player) dequeue() *Song {
	p.queueList = p.queueList[1:]
	return <-p.queue
}
