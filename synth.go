package herosynth

import (
    "container/list"
    //"fmt"
)

const defaultPolyphony = 32

type HeroSynth struct {
    sampleRate float64
    polyphony int
    runningVoices *list.List
    availableVoices *list.List
    voicesByKey map[int][]*Voice
    invNumRunningVoices float64
}

func CreateHeroSynth(sampleRate float64) *HeroSynth {
    synth := &HeroSynth{
        sampleRate: sampleRate,
        polyphony: defaultPolyphony,
        runningVoices: list.New(),
        availableVoices: list.New(),
        voicesByKey: make(map[int][]*Voice),
        invNumRunningVoices: 0.0}

    patch := &Patch{}

    for i := 0; i < defaultPolyphony; i++ {
        var voice *Voice = CreateVoice(synth, patch)
        voice.id = i
        synth.availableVoices.PushBack(voice)
        //fmt.Println(voice)
    }

    return synth
}

func (synth *HeroSynth) Render(out [][]float32) {

    for i := range out[0] {
        acum := [2]float64{0, 0}

        for e := synth.runningVoices.Front(); e != nil; e = e.Next() {
            voice := e.Value.(*Voice)
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

func (synth *HeroSynth) allocVoice() *Voice {
    //TODO add key as a parameter and update voices by key
    var voice *Voice = nil;

    if synth.availableVoices.Len() > 0 {
        e := synth.availableVoices.Front()
        voice = e.Value.(*Voice)
        synth.availableVoices.Remove(e)
    } else if synth.runningVoices.Len() > 0 {
        // TODO There are no available voices, kill the oldest one
        voice = synth.runningVoices.Front().Value.(*Voice)
    } else {
        return nil
    }

    synth.runningVoices.PushBack(voice)
    synth.invNumRunningVoices = 1.0 / float64(synth.runningVoices.Len())

    // TODO update patch

    return voice
}

func (synth *HeroSynth) noteOn(key uint, velocity float64) {
    voice := synth.allocVoice()
    if voice != nil {
        voice.noteOn(key, velocity)
    }
    //fmt.Println("---------------------------------------")
    //for e := synth.runningVoices.Front(); e != nil; e = e.Next() {
    //    voice := e.Value.(*Voice)
    //    fmt.Println(voice)
    //}
}

func (synth *HeroSynth) noteOff(key uint) {

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
