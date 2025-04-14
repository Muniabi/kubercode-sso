package values

import (
	"encoding/json"
	"testing"
)

func TestNewIsMentor(t *testing.T) {
	isMentor, err := NewIsMentor(true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if isMentor.IsMentor != true {
		t.Errorf("expected true, got %v", isMentor.IsMentor)
	}
}

func TestIsMentor_UnmarshalJSON(t *testing.T) {
	data := []byte(`{"is_mentor": true}`)
	var isMentor IsMentor
	err := json.Unmarshal(data, &isMentor)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if isMentor.IsMentor != true {
		t.Errorf("expected true, got %v", isMentor.IsMentor)
	}
}

func TestIsMentor_ToString(t *testing.T) {
	isMentor, err := NewIsMentor(true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if isMentor.ToString() != "true" {
		t.Errorf("expected 'true', got %v", isMentor.ToString())
	}
}

func TestIsMentor_GetIsMentor(t *testing.T) {
	isMentor, err := NewIsMentor(true)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if isMentor.GetIsMentor() != true {
		t.Errorf("expected true, got %v", isMentor.GetIsMentor())
	}
} 