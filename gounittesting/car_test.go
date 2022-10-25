package gounittesting

import "testing"

// Simple testing what different between Fatal and Error
func TestNew(t *testing.T) {
	c, err := New("Matt", 100)
	if err != nil {
		t.Fatal("got errors:", err)
	}

	if c == nil {
		t.Error("car should be nil")
	}
}
