package fn_utils

import (
	"fmt"
	"sync"
)

var globalGoRoutineContext sync.Map

var globalGoRoutineTemplate = "%d"

func Set(key string, value interface{}) {
	if load, ok := globalGoRoutineContext.Load(fmt.Sprintf(globalGoRoutineTemplate, GetGoRoutineId())); !ok {
		var goRoutine sync.Map
		globalGoRoutineContext.Store(fmt.Sprintf(globalGoRoutineTemplate, GetGoRoutineId()), &goRoutine)
		goRoutine.Store(key, value)
	} else {
		goRoutine := load.(*sync.Map)
		goRoutine.Store(key, value)
	}
}

func Get(key string) (interface{}, bool) {
	if load, ok := globalGoRoutineContext.Load(fmt.Sprintf(globalGoRoutineTemplate, GetGoRoutineId())); !ok {
		return nil, false
	} else {
		goRoutine := load.(*sync.Map)
		return goRoutine.Load(key)
	}
}

func Copy(resource uint64, target uint64) {
	if load, ok := globalGoRoutineContext.Load(fmt.Sprintf(globalGoRoutineTemplate, resource)); ok {
		goRoutine := load.(*sync.Map)
		var targetRoutine sync.Map
		goRoutine.Range(func(key, value interface{}) bool {
			targetRoutine.Store(key, value)
			return true
		})
		globalGoRoutineContext.Store(fmt.Sprintf(globalGoRoutineTemplate, target), &targetRoutine)
	}
}

func Clear() {
	globalGoRoutineContext.Delete(fmt.Sprintf(globalGoRoutineTemplate, GetGoRoutineId()))
}
