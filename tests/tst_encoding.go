package main

import "fmt"

import (
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/keybind"
)

func main() {
	X, _ := xgbutil.NewConn()

	keybind.Initialize(X)

	// should output "{"
	fmt.Println(string(keybind.LookupString(X, 1, 34)))
	// should output "["
	fmt.Println(string(keybind.LookupString(X, 0, 34)))

	fmt.Println("---------------------------------------")

	// should output "A"
	fmt.Println(string(keybind.LookupString(X, 1, 38)))
	// should output "a"
	fmt.Println(string(keybind.LookupString(X, 0, 38)))
}
