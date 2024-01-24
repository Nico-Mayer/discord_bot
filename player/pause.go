package player

import "fmt"

func (p *Player) Pause() error {
	if p.PlayerStatus != Playing {
		err := fmt.Errorf("CANT PAUSE, BECAUSE NO SONG IS PLAYED CURRENTLY")
		return err
	}

	if p.PlayerStatus == Paused {
		err := fmt.Errorf("SONG IS ALREADY PAUSED")
		return err
	}

	p.PlayerStatus = Paused
	p.currentStream.SetPaused(true)

	return nil
}
