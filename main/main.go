package main

import (
    "fmt"
    "code.google.com/p/portaudio-go/portaudio"
    "time"
    "herosynth"
)

const (
    sampleRate float64 = 44100
)

func main() {
    fmt.Println("Hero Synth")

    portaudio.Initialize()
    defer portaudio.Terminate()

    synth := herosynth.CreateHeroSynth(sampleRate)

    var stream, err = portaudio.OpenDefaultStream(0, 2, sampleRate, 0, func (out [][]float32) {
        synth.Render(out)
    })

    if err != nil {
        fmt.Println("Error creating the audio stream")
        panic(err)
    }

    defer stream.Close()

    stream.Start()

    for i := uint(20); i < 20 + 12; i++ {
        synth.Send(herosynth.NoteOn{i, 1.0})
        time.Sleep(500 * time.Millisecond)
    }

    stream.Stop()
}