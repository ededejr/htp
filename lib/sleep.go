package lib

import "time"

func Sleep(seconds int) {
	timer := time.NewTimer(time.Duration(seconds) * time.Second)
	<-timer.C
}
