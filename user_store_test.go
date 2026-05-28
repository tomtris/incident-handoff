package main

import (
	"context"
	"errors"
	"testing"
)

func TestGetByUsername(t *testing.T) {
	var seedUsers = []User{
		{ID: "u1", Username: "anh", Password: hashPassword("anh123"), Role: "engineer"},
		{ID: "u2", Username: "bernd", Password: hashPassword("bernd123"), Role: "engineer"},
		{ID: "u3", Username: "admin", Password: hashPassword("admin123"), Role: "admin"},
	}
	users := NewInMemoryUserStore(seedUsers)
	for _, each := range users.users {
		_, err := users.GetByUsername(context.Background(), each.Username)
		if err != nil {
			t.Fatalf("expect no error")
		}
	}

	_, err := users.GetByUsername(context.Background(), "not-exist-user")
	if errors.Is(err, ErrUserNotFound) == false {
		t.Fatalf("expect ErrUserNotFound")
	}
}
