package herosynth

import (
    "math"
    //"fmt"
)

var keyFreq [128]float64

func init() {
    for i := 0; i < 128; i++ {
        keyFreq[i] = 440.0 * math.Pow(2.0, (float64(i) - 69) / 12.0)
    }
}

type Voice struct {

    synth *HeroSynth

    id int

    key uint
    amplitude float64

    osc []Oscillator

    oscOut []float64

    env []Envelope

    envOut []float64
}

func CreateVoice(synth *HeroSynth, patch *Patch) *Voice {
    voice := &Voice{synth: synth}
    voice.UpdatePatch(patch)
    return voice
}

func (v *Voice) UpdatePatch(patch *Patch) {
    // TODO

    v.osc = make([]Oscillator, 1, 16)
    waveTable := NewSineWaveTable(DefaultWaveTableSize)
    v.osc[0] = MakeOscillator(v.synth.sampleRate, waveTable)

    v.oscOut = make([]float64, 1, 16)

    v.env = make([]Envelope, 1, 16)

    v.envOut = make([]float64, 1, 16)

    //v.UpdateKeyAndVelocity(32, 1.0)
}

func (v *Voice) noteOn(key uint, velocity float64) {
    v.UpdateKeyAndVelocity(key, velocity)

    // TODO reset envelopes, lfos, ...
}

func (v *Voice) UpdateKeyAndVelocity(key uint, velocity float64) {
    v.key = key
    v.amplitude = velocity
    freq := keyFreq[key & 0x7f]
    for i := range v.osc {
        osc := &v.osc[i]
        if !osc.fixedFreq {
            osc.baseFrequency = freq
            osc.UpdateFrequency()
        }
    }
}

func (v *Voice) NextSample() (out [2]float64) {

    // Generate envelope signals
    for i, env := range v.env {
        v.envOut[i] = env.NextSample()
    }

    // Generate oscillator signals and mix and pan carrier outputs
    activeOsc := 0
    for i := range v.osc {
        var osc *Oscillator = &v.osc[i]
        signal := osc.NextSample()
        signal *= v.envOut[osc.envIndex]
        v.oscOut[i] = signal

        if osc.mixAmp > 0.0 {
            signal *= osc.mixAmp * v.amplitude

            out[0] += signal * osc.panLeft
            out[1] += signal * osc.panRight

            activeOsc++
        }
        //fmt.Printf("%d: %f %f (%f, %f) %f\n", i, v.oscOut[i], signal, out[0], out[1], osc.tableOffset)
    }

    // Calculate phase modulation
    for i := range v.osc {
        var carrier *Oscillator = &v.osc[i]
        var phaseOffset float64 = 0
        // TODO add LFOs to the phase
        for j := range v.osc {
            var modulator *Oscillator = &v.osc[j]
            phaseOffset += v.oscOut[j] * modulator.phaseMod
        }
        carrier.ModulatePhase(phaseOffset)
    }

    return out
}
