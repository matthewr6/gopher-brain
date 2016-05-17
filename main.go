package main

import (
    "fmt"
    "time"
    "bufio"
    "os"
    "strings"
    // "strconv"
    // "reflect"
    "math/rand"

    "github.com/jteeuwen/keyboard/termbox"
    term "github.com/nsf/termbox-go"
)

/*
    & is "address of"
    * is "value at address"
*/

var running = true

// todo - clean up goroutines and whatnot
func main() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter state name to load state, or leave blank to create a new network:  ")
    fileName, _ := reader.ReadString('\n')
    fileName = strings.TrimSpace(fileName)

    start := time.Now()
    rand.Seed(time.Now().UTC().UnixNano())

    // this is the keyboard sensing stuff
    term.Init()
    kb := termbox.New()
    kb.Bind(func() {
        running = false
    }, "space")
    go KeyboardPoll(kb)
    
    var myNet *Network
    if fileName != "" {
        NETWORK_SIZE := [3]int{25, 25, 25}
        myNet = MakeNetwork(NETWORK_SIZE, false)
        myNet.Connect()
    } else {
        myNet = LoadState(fileName)
    }
    myNet.CreateSensor(1, 50, "", [3]int{25, 1, 1}, true, "a", kb)
    myNet.CreateSensor(1, 50, "", [3]int{1, 1, 1}, true, "s", kb)
    myNet.AnimateUntilDone(100)
    
    // uncomment this to test state saving/loading capabilities
    // myNet.SaveState("test")
    // loadedNet := LoadState("test")
    // fmt.Println(reflect.DeepEqual(loadedNet, myNet))

    elapsed := time.Since(start)
    term.Close()

    fmt.Print("\nSave state?  Enter a name if you wish to save the state:  ")
    fileName, _ = reader.ReadString('\n')
    fileName = strings.TrimSpace(fileName)
    if fileName != "" {
        myNet.SaveState(fileName)
    }

    fmt.Printf("Took %s\n", elapsed)
}