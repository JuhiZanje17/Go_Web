//https://blog.golang.org/defer-panic-and-recover
/*
1:A deferred function's arguments are evaluated when the defer statement is evaluated.
2:Deferred function calls are executed in Last In First Out order after the surrounding function returns.
3:Deferred functions may read and assign to the returning function's named return values.
*/

package main

import (
	"fmt"
)
//1
func a() {
    i := 0 
    defer fmt.Println(i)   
    i++
    //defer fmt.Println(i)
    return
}
//3:this fn will return 3
func c() (i int) {
    defer func() { i+=2 }()
    return 1
}

func main(){
    fmt.Println(c())
}



