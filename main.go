// _Timeouts_ are important for programs that connect to
// external resources or that otherwise need to bound
// execution time. Implementing timeouts in Go is easy and
// elegant thanks to channels and `select`.

package main

import (
	"errors"
	"fmt"
	"time"
)

func main() {
	err, response := TestTimeout(1, 10)
	if err != nil {
		fmt.Printf("%s\t(error)\n", err)
	} else {
		fmt.Printf("%s\n", response)
	}


	err, response = TestTimeout(10, 2)
	if err != nil {
		fmt.Printf("%s\t(error)\n", err)
	} else {
		fmt.Printf("%s\n", response)
	}


	err, response = TestTimeout(3, 3)
	if err != nil {
		fmt.Printf("%s\t(error)\n", err)
	} else {
		fmt.Printf("%s\n", response)
	}
}


func TestTimeout(sleepTime int, timeout int) (error, string) {
	myChannel := make(chan string)
	go func() {
		time.Sleep(time.Duration(sleepTime) * time.Second)
		myChannel <- fmt.Sprintf("Sleep: %d\tTimeout: %d\t- Complete", sleepTime, timeout)
	}()
	select {
	case response := <-myChannel:
		return nil, response
	case <-time.After(time.Duration(timeout) * time.Second):
		return errors.New(fmt.Sprintf("Sleep: %d\tTimeout: %d\t- Timed Out", sleepTime, timeout)), ""
	}
}