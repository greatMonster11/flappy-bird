package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/veandco/go-sdl2/mix"
)

const (
    soundPath = "res/sounds/"
)

type audio struct {}

func newAudio() (*audio, error) {
    if err := mix.OpenAudio(44100, mix.DEFAULT_FORMAT, 2, 4096); err != nil {
		return nil, fmt.Errorf("could not copy background: %v", err)
	}
	defer mix.CloseAudio()

    return &(audio{}),  nil
}

func playSound(action string) {
    soundURL := soundPath + action + ".wav"
    // fmt.Println(soundURL)

    data, err := ioutil.ReadFile(soundURL)
    if err != nil {
        log.Println(err)
    }

    // Load WAV from data (memory)
	chunk, err := mix.QuickLoadWAV(data)
	if err != nil {
		log.Println(err)
	}
	defer chunk.Free()

    // Play the chunk once.
    if _, err := chunk.Play(-1,0); err != nil {
        log.Println(err)
    }
}

func (a *audio) playJump() {
    playSound("jump")
}

