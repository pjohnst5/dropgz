package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Starting sleep")
	// Need to mimic CNI more then, didn't get an issue
	time.Sleep(30 * time.Second)
	fmt.Println("done")
}
