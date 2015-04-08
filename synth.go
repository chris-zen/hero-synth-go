package main

import (
    "math/rand"
    "container/list"
)

const defaultPolyphony = 32

type Voice struct {
    key int
    amplitude float32
}

type HeroSynth struct {
    sampleRate float64
    polyphony int
    runningVoices *list.List
    availableVoices *list.List
    voicesByKey map[int][]*Voice
    invNumRunningVoices float32
}

func CreateVoice() *Voice {
    voice := &Voice{}
    return voice
}

func (*Voice) NextSample() (out [2]float32) {
    out[0] = 2 * (rand.Float32() - 0.5)
    out[1] = 2 * (rand.Float32() - 0.5)
    return out
}

func CreateHeroSynth(sampleRate float64) *HeroSynth {
    synth := &HeroSynth{
        sampleRate: sampleRate,
        polyphony: defaultPolyphony,
        runningVoices: list.New(),
        availableVoices: list.New(),
        voicesByKey: make(map[int][]*Voice),
        invNumRunningVoices: 0.0}

    for i := 0; i < defaultPolyphony; i++ {
        synth.availableVoices.PushBack(CreateVoice())
    }

    return synth
}

func (synth *HeroSynth) Render(out [][]float32) {

    for i := range out[0] {
        acum := [2]float32{0, 0}

        for e := synth.runningVoices.Front(); e != nil; e = e.Next() {
            voice := e.Value.(*Voice)
            sample := voice.NextSample()
            acum[0] += sample[0]
            acum[1] += sample[1]
        }

        acum[0] *= synth.invNumRunningVoices
        acum[1] *= synth.invNumRunningVoices

        out[0][i] = acum[0]
        out[1][i] = acum[1]
    }
}

func (synth *HeroSynth) allocVoice() *Voice {
    var voice *Voice = nil;

    if synth.availableVoices.Len() > 0 {
        voice = synth.availableVoices.Front().Value.(*Voice)
    } else if synth.runningVoices.Len() > 0 {
        // TODO There are no available voices, kill the oldest one
        voice = synth.runningVoices.Front().Value.(*Voice)
    } else {
        return nil
    }

    // TODO update patch

    return voice
}

func (synth *HeroSynth) noteOn(key int, velocity float32) {
    voice := synth.allocVoice()
    if voice != nil {
        voice.key = key
        voice.amplitude = velocity
        synth.runningVoices.PushBack(voice)
        synth.invNumRunningVoices = 1.0 / float32(synth.runningVoices.Len())
    }
}

func (synth *HeroSynth) noteOff(key int) {

}

func (synth *HeroSynth) Send(event Event) {
    switch event.(type) {
        case NoteOn:
            noteOn := event.(NoteOn)
            synth.noteOn(noteOn.key, noteOn.velocity)
        case NoteOff:
            noteOff := event.(NoteOff)
            synth.noteOff(noteOff.key)

    }
}