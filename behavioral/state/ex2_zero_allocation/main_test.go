package main

import "testing"

func TestDocument_Workflow(t *testing.T) {
	// Scenario A (Happy path)
	doc := NewDocument()
	if err := doc.RequestReview(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if state := doc.GetStateName(); state != "Review" {
		t.Errorf("Expected state 'Review',got %s", state)
	}
	if err := doc.Approve(); err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	// Scenario B(Invalid Transition)
	doc = NewDocument()
	if err := doc.Approve(); err == nil {
		t.Error("Expected error, got nil")
	} else if err.Error() != "cannot approve a draft directly" {
		t.Errorf("Expected error 'cannot approve a draft directly', got %v", err)
	}
	doc.RequestReview()
	if err := doc.RequestReview(); err == nil {
		t.Error("Expected error, got nil")
	} else if err.Error() != "already in review" {
		t.Errorf("Expected error 'already in review', got %v", err)
	}
	doc.Approve()
	if err := doc.Approve(); err == nil {
		t.Error("Expected error, got nil")
	} else if err.Error() != "already published" {
		t.Errorf("Expected error 'already published',got %v", err)
	}

}
