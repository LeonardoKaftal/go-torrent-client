package bitfield

import "testing"

func TestBitfield_HavePiece(t *testing.T) {
	t.Log("Testing BitfieldHavePiece")
	input := Bitfield{0b00010000}
	result := input.HavePiece(3)
	if result != true {
		t.Errorf("Bitfield.HavePiece(1) returned %t, want %t", result, true)
	}
}

func TestBitfieldSetPiece(t *testing.T) {
	t.Log("Testing BitfieldSetPiece")
	input := Bitfield{0b00010000}
	input.SetPiece(6)
	if !input.HavePiece(6) {
		t.Error("Expected piece at 6 to be 1 but is 0")
	}
	input.SetPiece(12)
	// no crash
}
