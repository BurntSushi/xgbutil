package main

import "fmt"
import "github.com/BurntSushi/xgbutil"

func main() {
    X, _ := xgbutil.Dial("")

    heads, err := X.Heads()
    if err != nil {
        fmt.Printf("ERROR: %v\n", err)
    } else {
        for i, head := range heads {
            fmt.Printf("%d - %v\n", i, head)
        }
    }

    fmt.Println(X.WindowManager)
}

