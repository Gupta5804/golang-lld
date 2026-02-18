package main

import "testing"


func TestJob_Middleware_Chain(t *testing.T) {
	
	safeJob := Auth(Logging(Sanitize(BaseJob)),"admin")
	if output := safeJob("   data   ") ; output != "Processed: data"{
		t.Errorf("Expected 'Processed: data', got : %s", output)
	}

	restrictedJob := Auth(Logging(Sanitize(BaseJob)),"guest")
	
	if output := restrictedJob("hack") ; output != "Unauthorized" {
		t.Errorf("Expected 'Unauthorized', got : %s", output)
	}

}
