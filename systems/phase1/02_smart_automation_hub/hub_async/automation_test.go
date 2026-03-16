package hubasync_test

import (
	"sync"
	"testing"

	"github.com/Gupta5804/golang-lld/systems/phase1/02_smart_automation_hub/domain"
	// "github.com/Gupta5804/golang-lld/systems/phase1/02_smart_automation_hub/hub"
	hubasync "github.com/Gupta5804/golang-lld/systems/phase1/02_smart_automation_hub/hub_async"
)

func TestDispatcherhubasync(t *testing.T) {
	tests := []struct {
		name string
		// setup builds the topography and the pub-sub relationship
		setup          func() (*hubasync.DispatcherHub, map[string]domain.HubNode)
		triggerSensor  string
		triggerState   string
		expectedStates map[string]string // nodeID -> expectedState
	}{
		{
			name: "1-to-1 Mapping (Happy path)",
			setup: func() (*hubasync.DispatcherHub, map[string]domain.HubNode) {
				h := hubasync.NewDispatcherHub(1, 1)
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
			setup: func() (*hubasync.DispatcherHub, map[string]domain.HubNode) {
				h := hubasync.NewDispatcherHub(3, 1)
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
			setup: func() (*hubasync.DispatcherHub, map[string]domain.HubNode) {
				h := hubasync.NewDispatcherHub(2, 2)
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
			setup: func() (*hubasync.DispatcherHub, map[string]domain.HubNode) {
				h := hubasync.NewDispatcherHub(1, 1)
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
			myHub, nodesToInspect := tc.setup()

			myHub.Start()

			myHub.Trigger(tc.triggerSensor, tc.triggerState)

			myHub.Stop()

			for nodeName, expectedState := range tc.expectedStates {
				actualNode, exists := nodesToInspect[nodeName]
				if !exists {
					t.Fatalf("Test setup Error: Node %s not in assertion map", nodeName)
				}
				if actualState := actualNode.GetState(); actualState != expectedState {
					t.Errorf("Expected %s to be %s,got %s", nodeName, expectedState, actualState)
				}
			}
		})
	}
}

func TestDispatcherHub_GracefulShutdown(t *testing.T) {
	hub := hubasync.NewDispatcherHub(5, 500)
	hub.Start()

	targetDevice := domain.NewDevice("target_light")
	hub.Subscribe("sensor_1", targetDevice)

	var publisherWg sync.WaitGroup
	eventCount := 100

	publisherWg.Add(eventCount)
	for i := 0; i < eventCount; i++ {
		go func() {
			defer publisherWg.Done()
			hub.Trigger("sensor_1", "ON")
		}()
	}
	publisherWg.Wait()

	hub.Stop()

	if state := targetDevice.GetState() ; state != "ON" {
		t.Errorf("Expected device to be ON, got %s",state)
	}
}
