package main

import (
	"fmt"
	"log"
)

import (
	"github.com/BurntSushi/xgbutil"
)

func main() {
	_, err := xgbutil.Dial("")
	if err != nil {
		log.Fatalf("Could not connect to X: %v", err)
	}

	fmt.Println("Connected!")
}
