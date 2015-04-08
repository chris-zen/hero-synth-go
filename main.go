package main

import (
    "fmt"
    "code.google.com/p/portaudio-go/portaudio"
    "time"
)

const (
    sampleRate float64 = 44100
)

func main() {
    fmt.Println("Hero Synth")

    portaudio.Initialize()
    defer portaudio.Terminate()

    synth := CreateHeroSynth(sampleRate)

    var stream, err = portaudio.OpenDefaultStream(0, 2, sampleRate, 0, func (out [][]float32) {
        synth.Render(out)
    })

    if err != nil {
        fmt.Println("Error creating the audio stream")
        panic(err)
    }

    defer stream.Close()

    synth.Send(NoteOn{32, 1.0})
    stream.Start()
    time.Sleep(4 * time.Second)
    stream.Stop()
}