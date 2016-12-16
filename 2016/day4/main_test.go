package main

import (
	"testing"
	"strings"
)

func TestTrivialExample(t *testing.T) {
	text := `aaaaa-bbb-z-y-x-123[abxyz]
a-b-c-d-e-f-g-h-987[abcde]
not-a-real-room-404[oarel]
totally-real-room-200[decoy]
`
	data := gen(strings.Fields(text))
	program := &state{}
	fun := part1
	for fun != nil {
		fun = fun(program, data)
	}
	expected := 1514
	if program.answer != expected {
		t.Errorf("exprected %v, got %v; program state %v", expected,
			program.answer, program)
	}
}
