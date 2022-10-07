package lib

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func RunWorkers() {
	var wg sync.WaitGroup
	var count uint64

	for i := 0; i < 10; i++ {
		wg.Add(1)
		index := i
		go func() {
			id := SetInterval(func() {
				inc := func() {
					atomic.AddUint64(&count, 1)
				}
				worker(index, inc)
			}, 100)
			Sleep(10)
			ClearInterval(id)
			wg.Done()
		}()
	}

	wg.Wait()
	fmt.Println("Total requests:", count)
}

func worker(i int, inc func()) {
	// If the request below fails, this will
	// recover things and continue
	defer func() {
		if r := recover(); r != nil {
			fmt.Println(fmt.Sprintf("[%d] Recovered. Error:\n", i), r)
		}
	}()

	resp, err := makeMeasuredRequest(i, "https://google.com")
	if err != nil {
		panic(err)
	}
	fmt.Printf("[%d] %s\n", i, resp.Status)
	inc()
}
