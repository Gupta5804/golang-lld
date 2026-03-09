package hub_test

import (
	"fmt"
	"sync"
	"testing"

	"github.com/Gupta5804/golang-lld/systems/phase1/02_smart_automation_hub/domain"
	"github.com/Gupta5804/golang-lld/systems/phase1/02_smart_automation_hub/hub"
)

func TestAutomationHub_Routing(t *testing.T) {
	tests := []struct {
		name string
		// setup builds the topography and the pub-sub relationship
		setup          func() (*hub.AutomationHub, map[string]domain.HubNode)
		triggerSensor  string
		triggerState   string
		expectedStates map[string]string // nodeID -> expectedState
	}{
		{
			name: "1-to-1 Mapping (Happy path)",
			setup: func() (*hub.AutomationHub, map[string]domain.HubNode) {
				h := hub.NewAutomationHub()
				light := domain.NewDevice("light_1")

				h.Subscribe("sensor_1", light)

				return h, map[string]domain.HubNode{"light_1": light}
			},
			triggerSensor: "sensor_1",
			triggerState:  "ON",
			expectedStates: map[string]string{
				"light_1": "ON",
			},
		},
		{
			name: "Composite Cascade (Zone testing)",
			setup: func() (*hub.AutomationHub, map[string]domain.HubNode) {
				h := hub.NewAutomationHub()
				light1 := domain.NewDevice("light_1")
				light2 := domain.NewDevice("light_2")
				light3 := domain.NewDevice("light_3")

				zone := domain.NewZone("living_room")
				zone.AddNode(light1)
				zone.AddNode(light2)
				zone.AddNode(light3)
				h.Subscribe("sensor_1", zone)
				return h, map[string]domain.HubNode{
					"light_1": light1,
					"light_2": light2,
					"light_3": light3,
				}
			},
			triggerSensor: "sensor_1",
			triggerState:  "ON",
			expectedStates: map[string]string{
				"light_1": "ON",
				"light_2": "ON",
				"light_3": "ON",
			},
		},
		{
			name: "Many-to-Many test",
			setup: func() (*hub.AutomationHub, map[string]domain.HubNode) {
				h := hub.NewAutomationHub()
				light1 := domain.NewDevice("light_1")
				light2 := domain.NewDevice("light_2")
				h.Subscribe("sensor_1", light1)
				h.Subscribe("sensor_1", light2)

				return h, map[string]domain.HubNode{
					"light_1": light1,
					"light_2": light2,
				}
			},
			triggerSensor: "sensor_1",
			triggerState:  "ON",
			expectedStates: map[string]string{
				"light_1": "ON",
				"light_2": "ON",
			},
		},
		{
			name: "Isolation Principle",
			setup: func() (*hub.AutomationHub, map[string]domain.HubNode) {
				h := hub.NewAutomationHub()
				light1 := domain.NewDevice("light_1")
				light2 := domain.NewDevice("light_2")
				h.Subscribe("sensor_1", light1)
				return h, map[string]domain.HubNode{
					"light_1": light1,
					"light_2": light2,
				}
			},
			triggerSensor: "sensor_1",
			triggerState:  "ON",
			expectedStates: map[string]string{
				"light_1": "ON",
				"light_2": "OFF",
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// call tc.setup -> returns hub
			myHub, nodesToInspect := tc.setup()
			myHub.Trigger(tc.triggerSensor, tc.triggerState)

			for nodeName, expectedState := range tc.expectedStates {
				actualNode, exists := nodesToInspect[nodeName]
				if !exists {
					t.Fatalf("Test Setup Error: Node %s not found in assertion map", nodeName)
				}
				if actualState := actualNode.GetState(); actualState != expectedState {
					t.Errorf("Expected %s to be %s, got %s", nodeName, expectedState, actualState)
				}
			}
		})
	}
}

func TestAutomationHub_Concurrency(t *testing.T) {
	h := hub.NewAutomationHub()
	baselight := domain.NewDevice("base_light")
	h.Subscribe("sensor_1", baselight)

	var wg sync.WaitGroup
	chaosFactor := 100

	wg.Add(2 * chaosFactor)

	for i := 0; i < chaosFactor; i++ {
		// Readers (these will loop through the map)
		go func() {
			defer wg.Done()
			h.Trigger("sensor_1", "ON")
		}()

		go func(index int) {
			defer wg.Done()
			// device := domain.NewDevice("dynamic_light_" + string(index))
			device := domain.NewDevice(fmt.Sprintf("dynamic_light_%d",index))
			h.Subscribe("sensor_1", device)
		}(i)
	}
	wg.Wait()
}
