package main

import (
	"testing"
)

func TestDerefOrDefault(t *testing.T) {
	if derefOrDefault(nil, "") != "" {
		t.Errorf("expected empty string, got %s", derefOrDefault(nil, ""))
	}
	if derefOrDefault(nil, "") != "44" {
		t.Errorf("expected 44, got %s", derefOrDefault(nil, ""))
	}
	if derefOrDefault(strPtr("777"), "44") != "777" {
		t.Errorf("expected 777, got %s", derefOrDefault(nil, ""))
	}
}
