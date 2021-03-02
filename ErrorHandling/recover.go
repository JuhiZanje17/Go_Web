package main

import (
	"fmt"
)

func main(){
	f()
	fmt.Println("completed main f()")
}

func f(){

	defer func(){
		if r:=recover();r!=nil{
			fmt.Println("Recovered in f()",r)
		}
	}()
    
	fmt.Println("calling G")
	g(0)
	fmt.Println("completed main f()")
}    

func g(i int) {
    if i > 3 {
        fmt.Println("Panicking!")
        panic(i)
    }
    defer fmt.Println("Defer in g", i)
    fmt.Println("Printing in g", i)
    g(i + 1)       
}