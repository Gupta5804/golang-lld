package main

import (
	"testing"
)

// ---testing factory logic of NewDatabase---

func TestNewDatabase_Production(t *testing.T) {
	saver, err := NewDatabase("production")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	_, ok := saver.(*SQLDatabase)
	if !ok {
		t.Errorf("Expected *SQLDatabase, got %T", saver)
	}
}

func TestNewDatabase_Local(t *testing.T) {
	saver, err := NewDatabase("local")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	_, ok := saver.(*FileDatabase)
	if !ok {
		t.Errorf("Expected *FileDatabase, got %T", saver)
	}
}

func TestNewDatabase_Invalid(t *testing.T) {
	_, err := NewDatabase("invalid_env")
	if err == nil {
		t.Error("Expected error , but got nil")
	}
}

// -- testing the service ---

type mockUserSaver struct {
	SaveUserCalled bool
	CapturedName   string
}

func (m *mockUserSaver) SaveUser(name string) {
	m.SaveUserCalled = true
	m.CapturedName = name
}

func TestUserService_RegisterUser(t *testing.T) {
	mock := &mockUserSaver{}
	service := NewUserService(mock)


	service.RegisterUser("Alice")

	if !mock.SaveUserCalled {
		t.Error("Expected SaveUser to be called, but wasn't")
	}
	if mock.CapturedName != "Alice"{
		t.Errorf("Expected user 'Alice', got %s",mock.CapturedName)
	}
}
