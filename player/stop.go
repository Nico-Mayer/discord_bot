package player

func (p *Player) Stop() {
	p.voiceConn.Disconnect()
	p.queue = make(chan *Song, 100)
	p.QueueList = []string{}
	p.voiceConn = nil
	p.PlayerStatus = Resting
}
