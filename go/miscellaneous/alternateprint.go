//package main
//
//import "fmt"
//
//var num int = 1
//
//func out1(ch chan int) {
//	for {
//		ch <- num
//		fmt.Println("out1:", <-ch)
//		num += 1
//		if num > 100 {
//			return
//		}
//	}
//}
//
//func out2(ch chan int) {
//	for {
//		fmt.Println("out2:", <-ch)
//		num += 1
//		if num > 100 {
//			return
//		}
//		ch <- num
//	}
//}
//
//func main() {
//	ch := make(chan int)
//	go out1(ch)
//	go out2(ch)
//	for {
//	}
//}

package main

import (
	"fmt"
	"sync"
	"time"
)

var num int = 1
var locker = new(sync.Mutex)
var cond = sync.NewCond(locker)

func out1() {
	for {
		cond.L.Lock()
		cond.Wait()

		if num%2 == 0 {
			fmt.Println("out1:", num)
		}
		num++

		if num > 100 {
			cond.L.Unlock()
			return
		}
		cond.Signal()
		cond.L.Unlock()
	}
}

func out2() {
	for {
		cond.L.Lock()
		cond.Wait()
		if num%2 != 0 {
			fmt.Println("out2:", num)
		}
		num++
		if num > 100 {
			cond.L.Unlock()
			return
		}
		cond.Signal()
		cond.L.Unlock()
	}
}

func main() {
	go out1()
	go out2()
	//cond.Signal()
	time.Sleep(time.Second)
	fmt.Println("cond")
	cond.Signal()
	fmt.Println("done")
	for {
	}
}
