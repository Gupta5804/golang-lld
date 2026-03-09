package domain

import (
	"fmt"
	"strings"
	"sync"
)

// Component interface
type HubNode interface {
	TurnOn()
	TurnOff()
	GetState() string
	GetID() string
}

// Composite (Zone)
type Zone struct {
	nodes []HubNode
	id    string
	mu    sync.RWMutex
}

func (z *Zone) AddNode(node HubNode) {
	z.mu.Lock()
	defer z.mu.Unlock()
	z.nodes = append(z.nodes, node)
}
func (z *Zone) TurnOn() {
	z.mu.Lock()
	defer z.mu.Unlock()
	for _, node := range z.nodes {
		node.TurnOn()
	}
}
func (z *Zone) TurnOff() {
	z.mu.Lock()
	defer z.mu.Unlock()
	for _, node := range z.nodes {
		node.TurnOff()
	}
}
func (z *Zone) GetState() string {
	z.mu.RLock()
	defer z.mu.RUnlock()
	var builder strings.Builder
	for _, node := range z.nodes {
		builder.WriteString(fmt.Sprintf("Device-ID [%s]-> State[%s]\n", node.GetID(), node.GetState()))
	}
	return builder.String()
}
func (z *Zone) GetID() string {
	return z.id
}
func NewZone(id string) *Zone {
	return &Zone{
		id:    id,
		nodes: make([]HubNode, 0),
	}
}

// Leaf (Device)
type Device struct {
	id    string
	state string
	mu    sync.RWMutex
}

func (d *Device) TurnOn() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.state = "ON"
}
func (d *Device) TurnOff() {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.state = "OFF"
}
func (d *Device) GetState() string {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.state
}
func (d *Device) GetID() string {
	return d.id
}

func NewDevice(id string) *Device {
	return &Device{
		id:    id,
		state: "OFF",
	}
}
