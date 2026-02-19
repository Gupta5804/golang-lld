package main

import "testing"



func TestProxy_Protection_FreeUser(t *testing.T) {
	engineBootCount = 0 // reset spy

	proxy := &RendererProxy{
		userTier: "Free",
	}

	_, err := proxy.RenderVideo("vacation.mp4")

	if err == nil || err.Error() != "upgrade to premium" {
		t.Errorf("Expected 'upgrade to premium', got %v",err)
	}

	if engineBootCount != 0 {
		t.Errorf("Fatal: Heavy engine booted %d times for a free user",engineBootCount)
	}
}


func TestProxy_Virtual_PremiumUser(t *testing.T){
	engineBootCount = 0
	proxy := &RendererProxy{
		userTier: "Premium",
	}

	res1, _ := proxy.RenderVideo("wedding.mp4")
	if res1 != "Successfully rendered wedding.mp4" {
		t.Errorf("Render failed, got %s",res1)
	}

	// Should reuse the already booted engine
	proxy.RenderVideo("one.mp4")
	proxy.RenderVideo("two.mp4")
	
	if engineBootCount != 1 {
		t.Errorf("Expected exactly 1 engine boot, but got %d",engineBootCount)
	}
}
