package hub

import (
	"sync"

	"github.com/Gupta5804/golang-lld/systems/phase1/02_smart_automation_hub/domain"
)

// automationHub acts as a central registry and event router (mediator + observer)
type AutomationHub struct {
	subscriptions map[string][]domain.HubNode // sensor's ID -> subscribed nodes
	mu            sync.RWMutex
}

// initializer for the registry
func NewAutomationHub() *AutomationHub {
	return &AutomationHub{
		subscriptions: make(map[string][]domain.HubNode),
	}
}

func (a *AutomationHub) Subscribe(sensorID string, node domain.HubNode) {
	a.mu.Lock()
	defer a.mu.Unlock()

	a.subscriptions[sensorID] = append(a.subscriptions[sensorID], node)
}

func (a *AutomationHub) Trigger(sensorID, newstate string) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if subscribedNodes, ok := a.subscriptions[sensorID]; ok {
		for _, node := range subscribedNodes {
			switch newstate {
			case "ON":
				node.TurnOn()
			case "OFF":
				node.TurnOff()
			}
		}
	}
}
