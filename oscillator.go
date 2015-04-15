package herosynth

import (
    "math"
    //"fmt"
)

type Oscillator struct {
    enabled bool

    sampleRate float64

    envIndex uint

    mixAmp float64
    panLeft float64
    panRight float64

    waveTable WaveTable
    freqToTableIncr float64
    tableIncr float64
    initialPhase float64
    tableOffset float64
    amplitude float64       // Oscillator signal amplitude

    fixedFreq bool          // When it is true the baseFrequency doesn't change with noteOn
    baseFrequency float64   // Oscillator base frequency
    octaves int             // Number of octaves to shift from the baseFrequency
    semitones int           // Number of semitones to shift from the baseFrequency
    detune float64          // Fine shift from the basefrequency
    frequency float64       // Calculated from baseFrequency, and octaves, semitones and detune

    freqMod float64         // Amount of modulation to the carriers
    phaseMod float64        // Phase modulation calculated from frequency and freqMod
}

func MakeOscillator(sampleRate float64, waveTable WaveTable) Oscillator {
    osc := Oscillator{
        sampleRate: sampleRate,
        waveTable: waveTable,
        freqToTableIncr: float64(waveTable.Size()) / sampleRate}

    osc.Init()

    return osc
}

func (osc *Oscillator) Init() {
    osc.enabled = true
    osc.amplitude = 1.0
    osc.fixedFreq = false
    osc.baseFrequency = 440.0
    osc.mixAmp = 1.0
    osc.panLeft = 1.0
    osc.panRight = 1.0

    osc.UpdateFrequency()
    osc.ResetPhase()
}

func (osc *Oscillator) ResetPhase() {
    //osc.tableIncr = osc.frequency * osc.freqToTableIncr
    osc.tableOffset = (osc.initialPhase / (2 * math.Pi)) * float64(osc.waveTable.Size())
}

func (osc *Oscillator) UpdateFrequency() {
    pitchScale := math.Pow(2, (float64(osc.octaves * 1200.0 + osc.semitones) * 100.0 + osc.detune) / 1200.0)
    osc.frequency = osc.baseFrequency * pitchScale;
    if osc.frequency < 0.0 {
        osc.frequency = 0.0;
    }
    osc.tableIncr = osc.frequency * osc.freqToTableIncr
    osc.phaseMod = osc.freqMod * osc.tableIncr
}

func (osc *Oscillator) ModulatePhase(phaseOffset float64) {
    osc.tableOffset += phaseOffset
}

func (osc *Oscillator) NextSample() float64 {

    waveTableSize := float64(osc.waveTable.Size())
    if osc.tableOffset < 0 {
        osc.tableOffset += waveTableSize
    } else if osc.tableOffset >= waveTableSize {
        osc.tableOffset -= waveTableSize
    }

    var value float64 = 0.0
    if osc.enabled && osc.amplitude > 0 {
        value = osc.waveTable.Sample(osc.tableOffset) * osc.amplitude
    }

    osc.tableOffset += osc.tableIncr
    //fmt.Printf("%f, %f, %f\n", osc.tableOffset, osc.tableIncr, value)
    return value
}
