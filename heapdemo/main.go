package main

import (
	"fmt"
	"log"
	"time"

	"github.com/Dieterbe/profiletrigger/heap"
)

type Data struct {
	d    []byte
	prev *Data
}

func New(prev *Data, size int) *Data {
	return &Data{
		d:    make([]byte, size),
		prev: prev,
	}
}

// HungryAllocator allocates 1MB every second
func HungryAllocator() {
	var prev *Data
	for {
		n := New(prev, 1000000)
		prev = n
		time.Sleep(time.Duration(1) * time.Second)
	}
}

// LightAllocator allocates 100kB every second
func LightAllocator() {
	var prev *Data
	for {
		n := New(prev, 100000)
		prev = n
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func main() {
	fmt.Println("allocating 1100 kB every second. should hit the 10MB threshold within 10 seconds..")
	errors := make(chan error)
	trigger, _ := heap.New("/tmp/prof", 10000000, 60, time.Duration(1)*time.Second, errors)
	go trigger.Run()
	go HungryAllocator()
	go LightAllocator()
	for e := range errors {
		log.Fatal("profiletrigger heap saw error:", e)
	}
}
