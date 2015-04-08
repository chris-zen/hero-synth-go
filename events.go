package main

type Event interface {
}

type NoteOn struct {
    key int
    velocity float32
}

type NoteOff struct {
    key int
}