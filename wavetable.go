package herosynth

import (
    "math"
)

const DefaultWaveTableSize uint = 1 << 14

type WaveTable interface {
    Size() uint
    Sample(offset float64) float64
}

type SineWaveTable struct {
    data []float64
}

func NewSineWaveTable(size uint) *SineWaveTable {
    data := make([]float64, size)

    var step float64 = 1.0 / float64(size);
    for i := uint(0); i < size; i++ {
        var x = float64(i) * step;
        data[i] = math.Sin(x * 2.0 * math.Pi);
    }

    return &SineWaveTable{data}
}

func (wt *SineWaveTable) Size() uint {
    return uint(len(wt.data))
}

func (wt *SineWaveTable) Sample(offset float64) float64 {
    var pos uint = uint(math.Floor(offset));
    var nextPos uint = (pos + 1) % uint(len(wt.data));
    value := wt.data[pos];
    nextValue := wt.data[nextPos];
    nextValueDiff := nextValue - value;
    fraction := offset - float64(pos);
    return value + nextValueDiff * fraction;
}