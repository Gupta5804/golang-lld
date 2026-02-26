package main

import "testing"

func TestFunctionalPipeline(t *testing.T) {
	finalHandler := func(r Request) error {
		r["processed"] = "true"
		return nil
	}
	chain := BuildChain(finalHandler, RequireAuth, RequireJSON)

	validReq := map[string]string{
		"Authorization": "Bearer 123",
		"Content-Type":  "application/json",
	}
	authFailReq := map[string]string{
		"Authorization": "invalid",
		"Content-Type":  "application/json",
	}
	jsonFailReq := map[string]string{
		"Authorization": "Bearer 123",
		"Content-Type":  "text/plain",
	}
	if err := chain(validReq); err != nil {
		t.Fatalf("Expected error to be nil, got %v",err)
	}
	if processed,ok := validReq["processed"] ; !ok{
		t.Error("Expected processed value, got no value")
	} else if processed != "true"{
		t.Errorf("Expected processed value for valid request to be true, got %s",processed)
	}

	if err := chain(authFailReq); err == nil{
		t.Fatal("Expected error, got nil")
	} else if err.Error() != "unauthorized" {
		t.Errorf("Expected error 'unauthorized', got %v",err)
	}
	if _,ok := authFailReq["processed"];ok{
		t.Error("Expected no value for auth fail request processed")
	}

	if err := chain(jsonFailReq); err == nil{
		t.Fatal("Expected error, got nil")
	} else if err.Error() != "invalid content type" {
		t.Errorf("Expected error 'invalid content type', got %v", err)
	}

	if _,ok := jsonFailReq["processed"];ok{
		t.Error("Expected no value for json fail request processed")
	}
}
