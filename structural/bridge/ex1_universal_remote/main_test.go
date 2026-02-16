package main

import (
	"testing"
)

func TestUniversalRemote_TogglePower(t *testing.T) {
	tv := &TV{}
	remote := &UniversalRemote{device: tv}
	remote.TogglePower()
	if tv.IsEnable() != true {
		t.Errorf("Expected the power to be ON, but was OFF")
	}
	remote.TogglePower()
	if tv.IsEnable() != false {
		t.Errorf("Expected the power to be OFF , but was ON")
	}

}

func TestRadio_TogglePower(t *testing.T) {
	radio := &Radio{}
	remote := &UniversalRemote{device: radio}
	remote.TogglePower()
	if radio.IsEnable() != true {
		t.Errorf("Expected the power to be ON, but was OFF")
	}
	remote.TogglePower()
	if radio.IsEnable() != false {
		t.Errorf("Expected the power to be OFF , but was ON")
	}
}

func TestAdvancedRemote_Mute(t *testing.T){
	tv := &TV{Volume:50}
	remote := &AdvancedRemote{UniversalRemote{device:tv}}
	remote.Mute()
	if tv.Volume != 0 {
		t.Errorf("Expected volume to be 0, but got :%d", tv.Volume)
	}
}
