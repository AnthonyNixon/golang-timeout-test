// _Timeouts_ are important for programs that connect to
// external resources or that otherwise need to bound
// execution time. Implementing timeouts in Go is easy and
// elegant thanks to channels and `select`.

package main

import (
	"bytes"
	"errors"
	"fmt"
	"os/exec"
	"strings"
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

	response, err = LongRunningBashProcessTimeout()
	if err != nil {
		fmt.Printf("%s\t(error)\n", err)
	} else {
		fmt.Printf("%s\n", response)
	}

	response, err = ShortRunningBashProcessTimeout()
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

func LongRunningBashProcessTimeout() (string, error) {
	cmd := exec.Command("./endless_sleep.sh")

	// Use a bytes.Buffer to get the output
	var buf bytes.Buffer
	cmd.Stdout = &buf

	cmd.Start()

	// Use a channel to signal completion so we can use a select statement
	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	// Start a timer
	timeout := time.After(time.Second * 3)

	// The select statement allows us to execute based on which channel
	// we get a message from first.
	select {
	case <-timeout:
		// Timeout happened first, kill the process and print a message.
		cmd.Process.Kill()
		return buf.String(), errors.New(fmt.Sprintf("bash function (pid %s) timed out", strings.Replace(buf.String(), "\n", "", -1)))
	case err := <-done:
		// Command completed before timeout. Print output and error if it exists.
		return buf.String(), err
	}

	return buf.String(), nil
}

func ShortRunningBashProcessTimeout() (string, error) {
	cmd := exec.Command("./2_second_sleep.sh")

	// Use a bytes.Buffer to get the output
	var buf bytes.Buffer
	cmd.Stdout = &buf

	cmd.Start()

	// Use a channel to signal completion so we can use a select statement
	done := make(chan error)
	go func() { done <- cmd.Wait() }()

	// Start a timer
	timeout := time.After(time.Second * 3)

	// The select statement allows us to execute based on which channel
	// we get a message from first.
	select {
	case <-timeout:
		// Timeout happened first, kill the process and print a message.
		cmd.Process.Kill()
		return buf.String(), errors.New("bash function timed out")
	case err := <-done:
		// Command completed before timeout. Print output and error if it exists.
		return buf.String(), err
	}

	return buf.String(), nil
}