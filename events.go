package herosynth

type Event interface {
}

type NoteOn struct {
    Key uint
    Velocity float64
}

type NoteOff struct {
    Key uint
}