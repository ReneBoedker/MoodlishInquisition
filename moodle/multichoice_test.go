package moodle

import (
	"testing"
)

func TestBalanceGrades(t *testing.T) {
	answers := []*Answer{
		NewAnswer("true", 1),
		NewAnswer("true", 1),
		NewAnswer("false", 0),
		NewAnswer("false", -10),
		NewAnswer("false", -20),
	}

	mc := NewMultiChoice("", 1, answers)

	// Check without penalties...
	mc.BalanceGrades(false)
	expected := [...]float64{50, 50, 0, -10, -20}
	for i, v := range mc.answers {
		if v.grade != expected[i] {
			t.Errorf("Without penalties: Got score %.6f, but expected %.6f", v.grade, expected[i])
		}
	}

	// ...and with
	mc.BalanceGrades(true)
	expected = [...]float64{50, 50, -50, -50, -50}
	for i, v := range mc.answers {
		if v.grade != expected[i] {
			t.Errorf("With penalties: Got score %.6f, but expected %.6f", v.grade, expected[i])
		}
	}
}
