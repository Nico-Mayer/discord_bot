package player

import "fmt"

func (p *Player) Skip() error {
	if len(p.queue) == 0 {
		err := fmt.Errorf("THE QUEUE IS EMPTY. PLEASE ADD SONGS BEFORE PROCEEDING")
		return err
	}

	p.skipInterrupt <- true
	p.PlayerStatus = Resting
	go p.Play()
	return nil
}
