package main

type Device interface {
	Enable()
	Disable()
	IsEnable() bool
	SetVolume(int)
}

type TV struct {
	isON   bool
	Volume int
}

func (t *TV) Enable()           { t.isON = true }
func (t *TV) Disable()          { t.isON = false }
func (t *TV) IsEnable() bool    { return t.isON }
func (t *TV) SetVolume(vol int) { t.Volume = vol }

type Radio struct {
	isON   bool
	Volume int
}

func (r *Radio) Enable()           { r.isON = true }
func (r *Radio) Disable()          { r.isON = false }
func (r *Radio) IsEnable() bool    { return r.isON }
func (r *Radio) SetVolume(vol int) { r.Volume = vol }

type UniversalRemote struct {
	device Device
}

func (u *UniversalRemote) TogglePower() {
	if u.device.IsEnable() {
		u.device.Disable()
	} else {
		u.device.Enable()
	}
}

type AdvancedRemote struct {
	UniversalRemote
}
func (a *AdvancedRemote) Mute(){
	a.device.SetVolume(0)
}
