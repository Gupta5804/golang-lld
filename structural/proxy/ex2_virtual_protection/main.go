package main

import "errors"

type VideoRenderer interface { //contract
	RenderVideo(videoName string) (string, error)
}
var engineBootCount int // global counter to spy on startups

type HeavyRenderer struct{} // real subject

func NewHeavyRenderer() *HeavyRenderer {
	engineBootCount++
	return &HeavyRenderer{}
}

func (h *HeavyRenderer) RenderVideo(videoName string) (string, error) {
	return "Successfully rendered " + videoName, nil
}

type RendererProxy struct {
	renderer VideoRenderer
	userTier string
}

func (r *RendererProxy) RenderVideo(videoName string) (string, error) {
	
	if r.userTier != "Premium" {
		return "", errors.New("upgrade to premium")
	}
	
	
	if r.renderer == nil {
		r.renderer = NewHeavyRenderer()
	} 
	return r.renderer.RenderVideo(videoName)
	
	
}
