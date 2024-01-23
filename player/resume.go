package player

import "fmt"

func (p *Player) Resume() error {
	if p.PlayerStatus != Paused {
		err := fmt.Errorf("PLAYER IS NOT PAUSED, SO IT CANT BE RESUMED")
		return err
	}

	p.currentStream.SetPaused(false)
	p.PlayerStatus = Playing
	return nil
}
