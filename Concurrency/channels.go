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
	<-ch

}

func main() {

	ch:=make(chan bool,channelCapacity)

	numArray := makeRange(0, 100)

	start := time.Now()

	for i, _ := range numArray {

		ch<-true
		go apiCall(i,ch)

	}

	for i:=0;i<channelCapacity;i++{
		ch<-true
	}

	elapsed := time.Since(start)
	log.Printf("Time taken %s", elapsed)

}

//i had helped him//i was helping him
//i like to wait about 1 hour//
//it was like i do not perform

//like
//it like keeps me like top of
