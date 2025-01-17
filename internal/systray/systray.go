package systray

import (
	"github.com/getlantern/systray"
	"github.com/judyas/sound-player/internal/player"
	"github.com/sqweek/dialog"
	"log"
	"os"
)

type Tray interface {
	Run()
}

type AppTray struct {
	player player.Player
}

func NewAppTray(p player.Player) *AppTray {
	return &AppTray{player: p}
}

func (a *AppTray) Run() {
	systray.Run(a.onReady, a.onExit)
}

func (a *AppTray) onReady() {
	systray.SetIcon(loadIcon())
	systray.SetTooltip("Sound Player")
	systray.SetTooltip("Sound Player running")

	playItem := systray.AddMenuItem("Play", "Play sound")
	stopItem := systray.AddMenuItem("Stop", "Stop sound")
	changeItem := systray.AddMenuItem("Change sound", "Change what sound should be played")
	rainMelody := changeItem.AddSubMenuItem("Rain", "play rain.wav")
	silentMelody := changeItem.AddSubMenuItem("Silent", "play silent.wav")
	customMelody := changeItem.AddSubMenuItem("Custom", "play custom melody")
	exitItem := systray.AddMenuItem("Exit", "Exit application")

	go func() {
		for {
			select {
			case <-playItem.ClickedCh:
				a.player.Play()
			case <-stopItem.ClickedCh:
				a.player.Stop()
			case <-rainMelody.ClickedCh:
				a.setMelody("assets/rain.wav")
			case <-silentMelody.ClickedCh:
				a.setMelody("assets/silent.wav")
			case <-customMelody.ClickedCh:
				a.chooseCustomMelody()
			case <-exitItem.ClickedCh:
				systray.Quit()
			}
		}
	}()
}

func (a *AppTray) setMelody(path string) {
	a.player.SetFile(path)
	a.player.Play()
}

func (a *AppTray) chooseCustomMelody() {
	path, err := dialog.File().Filter("WAV files", "wav").Title("Select a WAV file").Load()
	if err != nil {
		log.Printf("choise file: %v", err)
		return
	}

	a.setMelody(path)
}

func (a *AppTray) onExit() {
	a.player.Stop()
}

func loadIcon() []byte {
	iconPath := "assets/tray-icon.ico"
	data, err := os.ReadFile(iconPath)
	if err != nil {
		log.Fatalf("icon load: %v", err)
	}
	return data
}
