package notifier

import (
	"sync"

	"github.com/Gupta5804/golang-lld/systems/phase1/03_stock_alerts/domain"
)

type DispatcherHub struct {
	commandChan     <-chan domain.Command
	wg          sync.WaitGroup
	workerCount int
}

func NewDispatcherHub(commandChan <-chan domain.Command, workers int) *DispatcherHub {
	return &DispatcherHub{
		commandChan:commandChan,
		workerCount: workers,
	}
}

func (d *DispatcherHub) Start() {
	for i:=0;i<d.workerCount;i++{
		d.wg.Add(1)
		go func(){
			defer d.wg.Done()
			for cmd := range d.commandChan{
				cmd.Execute()
			}
		}()
	}
}
func (d *DispatcherHub) Stop() {
	d.wg.Wait()
}
