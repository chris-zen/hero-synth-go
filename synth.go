package herosynth

import (
    //"fmt"
)

const numKeys uint = 128

type HeroSynth struct {
    sampleRate float64
    voices [numKeys]Voice
    runningVoices map[uint]bool
    availableVoices map[uint]bool
    invNumRunningVoices float64
}

func CreateHeroSynth(sampleRate float64) *HeroSynth {
    synth := &HeroSynth{
        sampleRate: sampleRate,
        runningVoices: make(map[uint]bool),
        availableVoices: make(map[uint]bool),
        invNumRunningVoices: 0.0}

    patch := &Patch{}

    for i := uint(0); i < numKeys; i++ {
        var voice *Voice = CreateVoice(synth, patch)
        voice.key = i
        voice.UpdateKeyAndVelocity(i, 1.0)
        synth.voices[i] = *voice
        synth.availableVoices[i] = true
        //fmt.Println(voice)
    }

    return synth
}

func (synth *HeroSynth) Render(out [][]float32) {

    for i := range out[0] {
        acum := [2]float64{0, 0}

        for i := range synth.runningVoices {
            voice := &synth.voices[i]
            sample := voice.NextSample()
            acum[0] += sample[0]
            acum[1] += sample[1]
        }

        acum[0] *= synth.invNumRunningVoices
        acum[1] *= synth.invNumRunningVoices

        out[0][i] = float32(acum[0])
        out[1][i] = float32(acum[1])
    }
}

func (synth *HeroSynth) noteOn(key uint, velocity float64) {
    voice := &synth.voices[key]
    if _, isRunningVoice := synth.runningVoices[key]; !isRunningVoice {
        synth.runningVoices[key] = true
        delete(synth.availableVoices, key)
        synth.invNumRunningVoices = 1 / float64(len(synth.runningVoices))
    }
    voice.noteOn(key, velocity)
}

func (synth *HeroSynth) noteOff(key uint) {
    if _, isRunningVoice := synth.runningVoices[key]; isRunningVoice {
        delete(synth.runningVoices, key)
        synth.availableVoices[key] = true
    }
}

func (synth *HeroSynth) Send(event Event) {
    switch event.(type) {
        case NoteOn:
            noteOn := event.(NoteOn)
            synth.noteOn(noteOn.Key, noteOn.Velocity)
        case NoteOff:
            noteOff := event.(NoteOff)
            synth.noteOff(noteOff.Key)

    }
}
