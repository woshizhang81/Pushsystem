package main

import "fmt"


func Producer(ch chan int) {
	for i := 1; i <= 100; i++ {
		ch <- i
		}
	close(ch)
}

func Consumer(id int, ch chan int, done chan bool) {
	for {
		value, ok := <-ch
		if ok {
			fmt.Printf("id: %d, recv: %d\n", id, value)
		} else {
			fmt.Printf("id: %d, closed\n", id)
			break
			}
		}
	done <- true
}

func main() {
	ch := make(chan int, 3)

	coNum := 5
	done := make(chan bool, 1)

	for i := 1; i <= coNum; i++ {
		go Consumer(i, ch, done)
	}
	go Producer(ch)
	for i := 1; i <= coNum; i++ {
		<-done
		fmt.Println("entrance")
	}
}
