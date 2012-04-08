package main

import (
	"fmt"
	"log"
)

import (
	"burntsushi.net/go/xgbutil"
)

func main() {
	_, err := xgbutil.Dial("")
	if err != nil {
		log.Fatalf("Could not connect to X: %v", err)
	}

	fmt.Println("Connected!")
}
