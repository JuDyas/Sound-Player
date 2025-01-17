package main

import (
	"github.com/judyas/sound-player/internal/player"
	"github.com/judyas/sound-player/internal/systray"
)

func main() {
	playerSound := player.NewSoundPlayer()
	tray := systray.NewAppTray(playerSound)
	playerSound.Play()
	tray.Run()

}
