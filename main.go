package main

import (
	"fmt"
	"time"
)

type Message struct {
	Body      string
	Timestamp [3]int
	sender    int
}

func event(pid int, counter [3]int) [3]int {
	counter[pid] += 1
	fmt.Printf("Event in process pid=%v. Counter=%v\n", pid, counter)
	return counter
}

// func max(x, y int) int {
// 	if x < y {
// 		return y
// 	}
// 	return x
// }

 func calcTimestamp(recvTimestamp, counter [3]int) [3]int {
	for index := 0; index < 3; index++ {
		if recvTimestamp[index]>counter[index]{
			counter[index]=recvTimestamp[index]
		}

	}
	return counter 
 }

func sendMessage(ch chan Message, pid int, counter [3]int) [3]int {
	counter[pid] += 1
	ch <- Message{"Test msg!!!", counter, pid}
	fmt.Printf("Message sent from pid=%v. Counter=%v\n", pid, counter)
	return counter

}

func receiveMessage(ch chan Message, pid int, counter [3]int) [3]int {
	message := <-ch
	counter = calcTimestamp(message.Timestamp, counter)
	counter[pid] += 1
	fmt.Printf("Message received at pid=%v. Counter=%v\n", pid, counter)
	return counter
}

func processOne(ch12, ch21 chan Message) {
	pid := 0
	counter := [3]int{0, 0, 0}
	counter = event(pid, counter)
	counter = sendMessage(ch12, pid, counter)
	counter = event(pid, counter)
	counter = receiveMessage(ch21, pid, counter)
	counter = event(pid, counter)

}

func processTwo(ch12, ch21, ch23, ch32 chan Message) {
	pid := 1
	counter := [3]int{0, 0, 0}
	
	counter = receiveMessage(ch12, pid, counter)
	counter = sendMessage(ch21, pid, counter)
	counter = sendMessage(ch23, pid, counter)
	counter = receiveMessage(ch32, pid, counter)

}

func processThree(ch23, ch32 chan Message) {
	pid := 2
	counter := [3]int{0, 0, 0}
	counter = receiveMessage(ch23, pid, counter)
	counter = sendMessage(ch32, pid, counter)

}

func main() {
	fmt.Println("Start Vetorial Clocks")
	oneTwo := make(chan Message, 100)
	twoOne := make(chan Message, 100)
	twoThree := make(chan Message, 100)
	threeTwo := make(chan Message, 100)

	go processOne(oneTwo, twoOne)
	go processTwo(oneTwo, twoOne, twoThree, threeTwo)
	go processThree(twoThree, threeTwo)

	time.Sleep(5 * time.Second)
}
