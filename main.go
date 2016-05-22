package main

import (
    "fmt"
    "time"
    "bufio"
    "os"
    "strings"
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

// var choice string
// fmt.Printf("Network has %v sensor(s).\n", len(myNet.Sensors))
// for _, sensor := range myNet.Sensors {
//     fmt.Printf("    %v\n", sensor.Name)
// }
// fmt.Println("\n    Add sensor? y/n")
// choice, _ = reader.ReadString('\n')
// if choice == "y" {

// }
// fmt.Println("    Remove sensor? y/n")
// choice, _ = reader.ReadString('\n')
// if choice == "y" {

// }

// fmt.Printf("Network has %v output(s).\n", len(myNet.Sensors))
// for _, output := range myNet.Outputs {
//     fmt.Printf("    %v\n", output.Name)
// }
// fmt.Println("\n    Add output? y/n")
// choice, _ = reader.ReadString('\n')
// if choice == "y" {

// }
// fmt.Println("    Remove output? y/n")
// choice, _ = reader.ReadString('\n')
// if choice == "y" {

// }

func main() {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print("Enter state name to load state, or leave blank to create a new network:  ")
    fileName, _ := reader.ReadString('\n')
    fileName = strings.TrimSpace(fileName)

    start := time.Now()
    rand.Seed(time.Now().UTC().UnixNano())

    
    var myNet *Network
    _, err := os.Stat(fmt.Sprintf("./state/%v_state.json", fileName))
    if fileName == "" || err != nil {
        NETWORK_SIZE := [3]int{25, 25, 25}
        myNet = MakeNetwork(NETWORK_SIZE, false)
        myNet.Connect()
        // myNet.CreateSensor("a", 1, 50, "", [3]int{25, 1, 1}, true, "a", kb)
        // myNet.CreateSensor("s", 1, 50, "", [3]int{1, 1, 1}, true, "s", kb)
        // myNet.CreateSensor("d", 1, 50, "", [3]int{12, 12, 1}, true, "d", kb)
        // myNet.CreateSensor("f", 1, 50, "", [3]int{12, 1, 12}, true, "f", kb)
        // myNet.CreateOutput("12, 1, 1", 1, 50,"", [3]int{12, 1, 1})
    } else {
        myNet = LoadState(fileName)
    }

    var choice string
    fmt.Printf("Network has %v sensor(s).\n", len(myNet.Sensors))
    for _, sensor := range myNet.Sensors {
        fmt.Printf("    %v\n", sensor.Name)
    }
    fmt.Println("\n    Add sensor? y/n")
    choice, _ = reader.ReadString('\n')
    if choice == "y" {

    }
    fmt.Println("    Remove sensor? y/n")
    choice, _ = reader.ReadString('\n')
    if choice == "y" {

    }

    fmt.Printf("Network has %v output(s).\n", len(myNet.Sensors))
    for _, output := range myNet.Outputs {
        fmt.Printf("    %v\n", output.Name)
    }
    fmt.Println("\n    Add output? y/n")
    choice, _ = reader.ReadString('\n')
    if choice == "y" {

    }
    fmt.Println("    Remove output? y/n")
    choice, _ = reader.ReadString('\n')
    if choice == "y" {

    }

    // this is the keyboard sensing stuff
    term.Init()
    term.SetCursor(0, 0)
    
    kb := termbox.New()
    kb.Bind(func() {
        running = false
        // if !netBuilt {
        //     os.Exit(1)
        // }
    }, "space")
    go KeyboardPoll(kb)
    myNet.BindKeyboard(kb)

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
    // myNet.CreateOutput("output", 1, 50,"", [3]int{12, 1, 1})
    // myNet.SaveState("test")
    // loadedNet := LoadState("test", nil)
    // fmt.Println(reflect.DeepEqual(loadedNet, myNet))
}