package main

import (
	"log"
	"time"

	"github.com/Dieterbe/profiletrigger/cpu"
)

var j int

func cpuUser(dur time.Duration, stopChan chan struct{}) {
	tick := time.NewTicker(time.Millisecond)
	//fmt.Print("every ms, sleep", dur, "   ")
	for {
		select {
		case <-tick.C:
			time.Sleep(dur)
		case <-stopChan:
			return
		// if we're not sleeping or stopping, just stay busy and consume cpu
		default:
			for i := 0; i < 10; i++ {
				j = i * 123
			}
		}
	}
}

func main() {
	errors := make(chan error)
	trigger, _ := cpu.New(".", 80, 60, time.Duration(1)*time.Second, time.Duration(2)*time.Second, errors)
	go trigger.Run()
	// gradually build up cpu usage.
	for i := 1000; i >= 0; i = i - 50 {
		stopChan := make(chan struct{})
		go cpuUser(time.Duration(i)*time.Microsecond, stopChan)
		time.Sleep(time.Second)
		stopChan <- struct{}{}
	}
	// end with just 100% cpu usage (1 core)
	go cpuUser(0, nil)

	for e := range errors {
		log.Fatal("profiletrigger cpu saw error:", e)
	}
}
