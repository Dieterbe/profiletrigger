package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Dieterbe/profiletrigger/heap"
)

type data struct {
	d    []byte
	prev *data
}

func newData(prev *data, size int) *data {
	return &data{
		d:    make([]byte, size),
		prev: prev,
	}
}

// HungryAllocator allocates 1MB every second
func HungryAllocator() {
	var prev *data
	for {
		n := newData(prev, 1000000)
		prev = n
		time.Sleep(time.Duration(1) * time.Second)
	}
}

// LightAllocator allocates 100kB every second
func LightAllocator() {
	var prev *data
	for {
		n := newData(prev, 100000)
		prev = n
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func main() {
	fmt.Println("allocating 1100 kB every second. should hit the 10MB threshold within 10 seconds.  look for a profile..")
	errors := make(chan error)
	trigger, _ := heap.New(".", 10000000, 60, time.Duration(1)*time.Second, errors)
	go trigger.Run()
	go HungryAllocator()
	go LightAllocator()
	for e := range errors {
		log.Fatal("profiletrigger heap saw error:", e)
	}
}
