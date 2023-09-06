package utilaio

import (
	"fmt"
	"sync"
)

type RoutineManager struct {
	max_workers  int
	workers      int
	manager_name string
	mutex        sync.Mutex
}

type Job interface {
	DoJob()
	QuitJob()
}

func (routine_manager *RoutineManager) incrementWorker() {
	routine_manager.mutex.Lock()
	routine_manager.workers += 1
	routine_manager.mutex.Unlock()
}

func (routine_manager *RoutineManager) decrementWorker() {
	routine_manager.mutex.Lock()
	routine_manager.workers -= 1
	routine_manager.mutex.Unlock()
}

func (routine_manager *RoutineManager) Execute(j Job) error {
	if routine_manager.workers < routine_manager.max_workers {
		go func() {
			routine_manager.incrementWorker()
			defer HandlePanic()
			j.DoJob()
			routine_manager.decrementWorker()
		}()
	} else {
		defer HandlePanic()
		j.QuitJob()
		return fmt.Errorf("all workers under %s routine manager are busy", routine_manager.manager_name)
	}

	return nil
}

func InitRoutineManager(manager_name string, max_workers int) *RoutineManager {
	routine_manager := &RoutineManager{
		manager_name: manager_name,
		max_workers:  max_workers,
		workers:      0,
	}
	return routine_manager
}
