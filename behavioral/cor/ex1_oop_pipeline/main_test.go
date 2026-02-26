package main

import "testing"

func TestPipeline_Execution(t *testing.T) {
	auth := &AuthHandler{}
	jsonHandler := &JSONHandler{}
	finalHandler := &FinalHandler{}

	auth.SetNext(jsonHandler).SetNext(finalHandler)
	validReq := &Request{
		Headers : map[string]string{
			"Authorization":"Bearer 123",
			"Content-Type":"application/json",
		},
		Body:"happy path",
	}
	authFailReq := &Request{
		Headers : map[string]string{
			"Authorization":"invalid",
			"Content-Type":"application/json",
		},
		Body:"Auth Failure",
	}
	jsonFailReq := &Request{
		Headers : map[string]string{
			"Authorization":"Bearer 123",
			"Content-Type":"text/plain",
		},
		Body:"JSON failure",
	}
	err1 :=auth.Execute(validReq)
	if err1 != nil {
		t.Fatalf("Expected err1 to be nil, got %v",err1)
	}
	if !validReq.IsProcessed{
		t.Error("IsProcessed for valid request is false")
	}
	err2 := auth.Execute(authFailReq)
	if err2 == nil {
		t.Fatal("Expected error, got nil")
	}
	if err2.Error() != "unauthorized"{
		t.Errorf("Expected error 'unauthorized', got %v",err2)
	}
	if authFailReq.IsProcessed{
		t.Error("IsProcessed for Auth fail request is true")
	}
	err3 := auth.Execute(jsonFailReq)
	if err3 == nil {
		t.Fatal("Expected error, got nil")
	}
	if err3.Error() != "invalid content type"{
		t.Errorf("Expected error 'invalid content type', got %v",err3)
	}
	if jsonFailReq.IsProcessed{
		t.Error("IsProcessed for json fail request is true")
	}
}
