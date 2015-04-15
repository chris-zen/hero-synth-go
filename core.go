package herosynth

type Synth interface {
    Render([][]float32)
}

func panAmp(pan float64) (leftAmp, rightAmp float64) {
    if pan > 0.0 {
        leftAmp = 1.0 - pan
        rightAmp = 1.0
    } else {
        leftAmp = 1.0
        rightAmp = 1.0 - pan
    }
    return leftAmp, rightAmp
}