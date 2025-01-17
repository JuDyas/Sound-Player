package player

import (
	"github.com/gopxl/beep/v2"
	"github.com/gopxl/beep/v2/speaker"
	"github.com/gopxl/beep/v2/wav"
	"os"
	"time"
)

type Player interface {
	SetFile(path string)
	Play()
	Stop()
}

type SoundPlayer struct {
	filePath string
	playing  bool
	stopChan chan bool
}

func NewSoundPlayer() *SoundPlayer {
	return &SoundPlayer{
		filePath: "assets/rain.wav",
		stopChan: make(chan bool),
	}
}

func (p *SoundPlayer) SetFile(path string) {
	p.filePath = path
}

func (p *SoundPlayer) Play() {
	if p.playing {
		p.Stop()
	}
	p.playing = true
	go func() {
		f, err := os.Open(p.filePath)
		if err != nil {
			p.playing = false
			panic(err)
			return
		}

		streamer, format, err := wav.Decode(f)
		if err != nil {
			p.playing = false
			panic(err)
			return
		}
		defer streamer.Close()

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
		loopBeep, err := beep.Loop2(streamer)
		if err != nil {
			p.playing = false
			panic(err)
			return
		}

		speaker.Play(beep.Seq(loopBeep, beep.Callback(func() {
			p.playing = false
		})))

		<-p.stopChan
		speaker.Clear()
	}()
}

func (p *SoundPlayer) Stop() {
	if p.playing {
		p.stopChan <- true
		p.playing = false
	}
}
