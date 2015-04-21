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

    const C1 = 24
    const C2 = 36
    const C3 = 48
    const C4 = 60
    for i := uint(0); i < 4; i++ {
        key := C3 + i * 12
        synth.Send(herosynth.NoteOn{key + 0, 1.0})
        synth.Send(herosynth.NoteOn{key + 4, 0.9})
        synth.Send(herosynth.NoteOn{key + 7, 0.8})
        time.Sleep(1000 * time.Millisecond)
        synth.Send(herosynth.NoteOff{key + 0})
        synth.Send(herosynth.NoteOff{key + 4})
        synth.Send(herosynth.NoteOff{key + 7})
    }

    stream.Stop()
}