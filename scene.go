package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	bg    *sdl.Texture
	bird  *bird
	pipes *pipes
	score *score
	audio *audio
}

func newScene(r *sdl.Renderer) (*scene, error) {
	bg, err := img.LoadTexture(r, "res/imgs/background.png")
	if err != nil {
		return nil, fmt.Errorf("could not load background image: %v", err)
	}

	bird, err := newBird(r)
	if err != nil {
		return nil, err
	}

	pipes, err := newPipes(r)
	if err != nil {
		return nil, err
	}

	score, err := newScore()
	if err != nil {
		return nil, err
	}

    audio, err := newAudio()
	if err != nil {
		return nil, err
	}

    return &scene{bg: bg, bird: bird, pipes: pipes, score: score, audio: audio}, nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) <-chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		ping := time.Tick(1000 * time.Millisecond)
		for {
			select {
			case e := <-events:
				if done := s.handleEvent(e); done {
					return
				}
				// log.Printf("event: %T", e)
            case <- ping:
                s.score.update()
			case <-tick:
				s.update()
				if s.bird.isDead() {
					drawTitle(r, "Game over")
					time.Sleep(time.Second)
					s.restart()
				}
				if err := s.paint(r); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (s *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {
	case *sdl.QuitEvent:
		return true
	case *sdl.MouseButtonEvent:
        s.audio.playJump()
		s.bird.jump()
	case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.TouchFingerEvent:
		// do nothing
	default:
		log.Printf("unknow event: %T", event)
	}
	return false
}

func (s *scene) update() {
	s.bird.update()
	s.pipes.update()
	s.pipes.touch(s.bird)
}

func (s *scene) restart() {
	s.bird.restart()
	s.pipes.restart()
	s.score.restart()
}

func (s *scene) paint(r *sdl.Renderer) error {
	r.Clear()

	if err := r.Copy(s.bg, nil, nil); err != nil {
		return fmt.Errorf("could not copy background: %v", err)
	}

	if err := s.bird.paint(r); err != nil {
		return err
	}

	if err := s.pipes.paint(r); err != nil {
		return err
	}

	if err := s.score.paint(r); err != nil {
		return err
	}

	defer r.Present()
	return nil
}

func (s *scene) destroy() {
	s.bg.Destroy()
	s.bird.destroy()
	s.pipes.destroy()
	s.score.destroy()
}
