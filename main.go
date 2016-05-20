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
    term.SetCursor(0, 0)
    
    netBuilt := false
    
    kb := termbox.New()
    kb.Bind(func() {
        running = false
        if !netBuilt {
            os.Exit(1)
        }
    }, "space")
    go KeyboardPoll(kb)
    
    var myNet *Network
    _, err := os.Stat(fmt.Sprintf("./state/%v_state.json", fileName))
    if fileName == "" || err != nil {
        NETWORK_SIZE := [3]int{25, 25, 25}
        myNet = MakeNetwork(NETWORK_SIZE, false)
        myNet.Connect()
        
        myNet.CreateSensor("aa", 1, 50, "", [3]int{25, 1, 1}, true, "a", kb)
        myNet.CreateSensor("bb", 1, 50, "", [3]int{1, 1, 1}, true, "b", kb)
    } else {
        myNet = LoadState(fileName, kb)
    }
    netBuilt = true
    myNet.AnimateUntilDone(100)
    

    elapsed := time.Since(start)
    term.Close()

    fmt.Print("\nSave state?  Enter a name if you wish to save the state:  ")
    fileName, _ = reader.ReadString('\n')
    fileName = strings.TrimSpace(fileName)
    if fileName != "" {
        myNet.SaveState(fileName)
    }

    fmt.Printf("Took %s\n", elapsed)

    // this section is to test state saving/loading capabilities
    // NETWORK_SIZE := [3]int{25, 25, 25}
    // myNet := MakeNetwork(NETWORK_SIZE, false)
    // myNet.Connect()

    // myNet.CreateSensor("aa", 1, 50, "", [3]int{24, 0, 0}, true, "a", nil)
    // myNet.CreateSensor("bb", 1, 50, "", [3]int{0, 0, 0}, true, "b", nil)
    // myNet.SaveState("test")
    // loadedNet := LoadState("test", nil)
    // fmt.Println(reflect.DeepEqual(loadedNet, myNet))
}