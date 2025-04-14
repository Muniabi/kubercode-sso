package values

import (
	"encoding/json"
	"testing"
)

func TestNewIsCompany(t *testing.T) {
	isCompany, err := NewIsCompany(true)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if isCompany.IsCompany != true {
		t.Errorf("expected true, got %v", isCompany.IsCompany)
	}
}

func TestUnmarshalJSON_Valid(t *testing.T) {
	data := []byte(`{"is_company": true}`)
	var isCompany IsCompany
	err := json.Unmarshal(data, &isCompany)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if isCompany.IsCompany != true {
		t.Errorf("expected true, got %v", isCompany.IsCompany)
	}
}

func TestUnmarshalJSON_Invalid(t *testing.T) {
	data := []byte(`{"is_company": "not_a_bool"}`)
	var isCompany IsCompany
	err := json.Unmarshal(data, &isCompany)
	if err == nil {
		t.Errorf("expected error for invalid JSON, got nil")
	}
}

func TestIsCompany_ToString(t *testing.T) {
	isCompany, err := NewIsCompany(true)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if isCompany.ToString() != "true" {
		t.Errorf("expected 'true', got %v", isCompany.ToString())
	}
}

func TestIsCompany_GetIsCompany(t *testing.T) {
	isCompany, err := NewIsCompany(true)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if isCompany.GetIsCompany() != true {
		t.Errorf("expected true, got %v", isCompany.GetIsCompany())
	}
}
