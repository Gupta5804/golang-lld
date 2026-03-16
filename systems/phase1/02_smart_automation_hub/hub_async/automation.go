package hubasync

import (
	"hash/fnv"
	"sync"

	"github.com/Gupta5804/golang-lld/systems/phase1/02_smart_automation_hub/domain"
)

type Event struct {
	SensorID string
	NewState string
}

type DispatcherHub struct {
	// the registry (state)
	subscriptions map[string][]domain.HubNode
	mu            sync.RWMutex

	// Async Engine (Behavior)
	workerCount    int
	workerChannels []chan Event //dedicated channel for each worker
	wg             sync.WaitGroup
}

func NewDispatcherHub(workerCount, bufferSize int) *DispatcherHub {
	workerChannels := make([]chan Event, workerCount)
	for i := 0; i < workerCount; i++ {
		workerChannels[i] = make(chan Event, bufferSize)
	}
	return &DispatcherHub{
		subscriptions:  make(map[string][]domain.HubNode),
		workerCount:    workerCount,
		workerChannels: workerChannels,
	}
}

// Start boots up the worker pool. should be called once at system startup
func (d *DispatcherHub) Start() {
	for i := 0; i < d.workerCount; i++ {
		d.wg.Add(1)
		go func(ch chan Event) {
			defer d.wg.Done()
			for event := range ch {
				d.mu.RLock()
				nodes, exists := d.subscriptions[event.SensorID]
				d.mu.RUnlock()
				if !exists {
					continue
				}
				for _, node := range nodes {
					switch event.NewState {
					case "ON":
						node.TurnOn()
					case "OFF":
						node.TurnOff()
					}
				}
			}
		}(d.workerChannels[i])
	}
}

// Initiates a graceful Shutdown
func (d *DispatcherHub) Stop() {
	for i := 0; i < d.workerCount; i++ {
		close(d.workerChannels[i])
	}
	d.wg.Wait()
}

func (d *DispatcherHub) Subscribe(sensorID string, node domain.HubNode) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.subscriptions[sensorID] = append(d.subscriptions[sensorID], node)
}

// Trigger no longer does the work, it acts as a "producer".
func (d *DispatcherHub) Trigger(sensorID string, newState string) {
	partition := hash(sensorID) % d.workerCount
	event := Event{SensorID: sensorID, NewState: newState}

	d.workerChannels[partition] <- event
}

// Helper function to hash a string to an integer
func hash(s string) int {
	h := fnv.New32a()
	h.Write([]byte(s))
	return int(h.Sum32())
}
