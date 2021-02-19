package main

import (
	"log"
	"time"
)

const channelCapacity=10

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func apiCall(i int,ch chan bool) {
	log.Println("API call for", i, "started")
	time.Sleep(100 * time.Millisecond)
	ch<-true
}

func main() {

	ch:=make(chan bool,channelCapacity)

	numArray := makeRange(0, 100)

	start := time.Now()

	for i, _ := range numArray {
		go apiCall(i,ch)		
	}

	for i:=0;i<channelCapacity;i++{
		<-ch
	}

	elapsed := time.Since(start)
	log.Printf("Time taken %s", elapsed)
}