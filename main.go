package main

import (
    "os"
    "fmt"
    "time"
    "bufio"
    "strings"
    "strconv"
    "math/rand"
)

/*
    & is "address of"
    * is "value at address"
*/

var running = true
var directory = "."

func main() {
    reader := bufio.NewReader(os.Stdin)
    fileName := Prompt("Enter state name to load state, or leave blank to create a new network:  ", reader)

    rand.Seed(time.Now().UTC().UnixNano())

    var myNet *Network
    _, err := os.Stat(fmt.Sprintf("./state/%v_state.json", fileName))
    if fileName == "" || err != nil {
        NETWORK_SIZE := [3]int{12, 25, 25}
        myNet = MakeNetwork(NETWORK_SIZE, false)
        myNet.Connect()
        myNet.Mirror()
        myNet.ConnectHemispheres()
    } else {
        myNet = LoadState(fileName)
    }
    var mode string
    tracker := make(map[string]bool)
    fmt.Println("Currently has the following sensors:")
    for name := range myNet.Sensors {
        baseName := strings.Split(name, "-")[0]
        if !tracker[baseName] {
            fmt.Println(baseName)
        }
        tracker[baseName] = true
    }
    mode = Prompt("Add/modify the custom things? [y/n]  ", reader)
    if mode == "y" {
        fmt.Println("WARNING!  Sensors and outputs will not save properly!")
        myNet.ClearIO() // is this needed
        // let's pretend the front x/z plane (y = 1) is "front" with left being x = 25
        // maybe you should only create sensors, and specify # of corresponding outputs - and then the createSensor generates the outputs automatically
        myNet.CreateSensor("eye", 1, 9, "y", [3]int{8, 0, 12}, 2, func(nodes []*Node, influences map[string]*Output) {
            for _, node := range nodes {
                node.Value = 1
            }
        })
    }

    framesInput := Prompt("Enter number of frames, or leave blank to run until manually stopped:  ", reader)
    frames, err := strconv.Atoi(framesInput)
    if err != nil {
        frames = 0
    }

    directory = Prompt("Enter directory to save frames and state to:  ", reader)
    if directory == "" {
        directory = "."
    }
    if directory[len(directory)-1] == '/' {
        directory = directory[0:len(directory)-1]
    }

    if frames == 0 {
        myNet.AnimateUntilDone()
    } else {
        myNet.GenerateAnim(frames)
    }

    fmt.Print("\nSave state?  Enter a name if you wish to save the state:  ")
    fileName, _ = reader.ReadString('\n')
    fileName = strings.TrimSpace(fileName)
    if fileName != "" {
        myNet.SaveState(fileName)
    }
}