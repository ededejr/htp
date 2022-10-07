package lib

import (
	"fmt"
	"sync"
	"time"
)

type intervalEntry struct {
	Stop func()
}

var (
	intervalMap      = make(map[string]intervalEntry)
	intervalMapMutex = sync.RWMutex{}
)

func SetInterval(fn func(), milliseconds int) string {
	ticker := time.NewTicker(time.Duration(milliseconds) * time.Millisecond)
	done := make(chan bool)

	go func() {
		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				fn()
			}
		}
	}()

	id := fmt.Sprintf("%d", time.Now().UnixNano())
	intervalMapMutex.Lock()
	intervalMap[id] = intervalEntry{
		Stop: func() {
			ticker.Stop()
			done <- true
		},
	}
	intervalMapMutex.Unlock()

	return id
}

func ClearInterval(id string) {
	intervalMapMutex.Lock()
	if entry, ok := intervalMap[id]; ok {
		entry.Stop()
		delete(intervalMap, id)
	}
	intervalMapMutex.Unlock()
}
