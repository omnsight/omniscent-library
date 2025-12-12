package helpers

import (
	"encoding/json"
	"testing"

	"github.com/omnsight/omniscent-library/gen/model/v1"
)

func TestJsonUnmarshalPerson(t *testing.T) {
	jsonData := `{"name": "Alice", "role": "Engineer", "nationality": "US", "birth_date": 1234567890}`

	var p model.Person
	err := json.Unmarshal([]byte(jsonData), &p)
	if err != nil {
		t.Fatalf("Failed to unmarshal Person: %v", err)
	}

	if p.Name != "Alice" {
		t.Errorf("Expected Name to be Alice, got %s", p.Name)
	}
	if p.Role != "Engineer" {
		t.Errorf("Expected Role to be Engineer, got %s", p.Role)
	}
	if p.Nationality != "US" {
		t.Errorf("Expected Nationality to be US, got %s", p.Nationality)
	}
	if p.BirthDate != 1234567890 {
		t.Errorf("Expected BirthDate to be 1234567890, got %d", p.BirthDate)
	}
}

func TestJsonUnmarshalRelation(t *testing.T) {
	jsonData := `{"_from": "persons/123", "_to": "organizations/456", "roles": ["founder"], "confidence": 95}`

	var r model.Relation
	err := json.Unmarshal([]byte(jsonData), &r)
	if err != nil {
		t.Fatalf("Failed to unmarshal Relation: %v", err)
	}

	if r.From != "persons/123" {
		t.Errorf("Expected From to be persons/123, got %s", r.From)
	}
	if r.To != "organizations/456" {
		t.Errorf("Expected To to be organizations/456, got %s", r.To)
	}
	if len(r.Roles) != 1 || r.Roles[0] != "founder" {
		t.Errorf("Expected Roles to contain 'founder', got %v", r.Roles)
	}
	if r.Confidence != 95 {
		t.Errorf("Expected Confidence to be 95, got %d", r.Confidence)
	}
}

func TestJsonUnmarshalEvent(t *testing.T) {
	jsonData := `{"title": "Conference", "description": "Tech conference", "happened_at": 1678886400}`

	var e model.Event
	err := json.Unmarshal([]byte(jsonData), &e)
	if err != nil {
		t.Fatalf("Failed to unmarshal Event: %v", err)
	}

	if e.Title != "Conference" {
		t.Errorf("Expected Title to be Conference, got %s", e.Title)
	}
	if e.Description != "Tech conference" {
		t.Errorf("Expected Description to be Tech conference, got %s", e.Description)
	}
	if e.HappenedAt != 1678886400 {
		t.Errorf("Expected HappenedAt to be 1678886400, got %d", e.HappenedAt)
	}
}
