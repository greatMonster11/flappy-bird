package main

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

const (
	fontPath = "./res/fonts/Flappy.ttf"
	fontSize = 6
)

type score struct {
	mu       sync.RWMutex
	textures []*sdl.Texture

	x, y  int32
	w, h  int32
	value  int64
}

func newScore() (*score, error) {
	var textures []*sdl.Texture

    return &score{textures: textures, y: 500, x: 10, value: 0}, nil
}

func (s *score) update() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.value += 1
}

func (s *score) restart() {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.value = 0
}

func (score *score) paint(r *sdl.Renderer) error {
	f, err := ttf.OpenFont(fontPath, fontSize)
	if err != nil {
		return fmt.Errorf("could not load font: %v", err)
	}
	defer f.Close()

	c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	s, err := f.RenderUTF8Solid(strconv.FormatInt(score.value, 10), c)
	if err != nil {
		return fmt.Errorf("could not render title: %v", err)
	}
	defer s.Free()

	t, err := r.CreateTextureFromSurface(s)
	if err != nil {
		return fmt.Errorf("could not create texture: %v", err)
	}
	defer t.Destroy()

	if err := r.Copy(t, nil, &sdl.Rect{X: 500, Y: 20, W: 50, H: 50}); err != nil {
	// if err := r.Copy(t, nil, nil); err != nil {
		return fmt.Errorf("could not copy texture: %v", err)
	}

	r.Present()

	return nil
}

func (s *score) destroy() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, t := range s.textures {
		t.Destroy()
	}
}


