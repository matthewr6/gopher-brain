package brain

import (
    "os"
    "fmt"
    "time"
    "bufio"
    "strings"
    "math/rand"
)

/*
    & is "address of"
    * is "value at address"
*/

var directory = "."

type SensorConstructor struct {
    Name string
    R int
    Count int
    Plane string
    Center [3]int
    OutputCount int
    InputFunc func([]*Node, map[string]*Output)
}

func Brain(NETWORK_SIZE [3]int, CONSTRUCTORS []SensorConstructor) *Network {
    reader := bufio.NewReader(os.Stdin)
    // todo - directory name to load state from
    fileName := Prompt("Enter state name to load state, or leave blank to create a new network:  ", reader)

    rand.Seed(time.Now().UTC().UnixNano())

    var myNet *Network
    _, err := os.Stat(fmt.Sprintf("./state/%v_state.json", fileName))
    if fileName == "" || err != nil {
        myNet = MakeNetwork(NETWORK_SIZE, false)
        myNet.Connect()
        myNet.Mirror()
        myNet.ConnectHemispheres()
    } else {
        myNet = LoadState(fileName)
    }
    tracker := make(map[string]bool)
    fmt.Printf("Currently has %v sensors and %v outputs.\n", len(myNet.Sensors), len(myNet.Outputs))
    if len(myNet.Sensors) > 0 {
        fmt.Println("Sensor names:")
        for name := range myNet.Sensors {
            baseName := strings.Split(name, "-")[0]
            if !tracker[baseName] {
                fmt.Println(baseName)
            }
            tracker[baseName] = true
        }
    }
    for _, constructor := range CONSTRUCTORS {
        myNet.CreateSensor(constructor.Name, constructor.R, constructor.Count, constructor.Plane, constructor.Center, constructor.OutputCount, constructor.InputFunc)
    }
    myNet.PruneUnusedSensors()
    fmt.Printf("Now has %v sensors and %v outputs.\n", len(myNet.Sensors), len(myNet.Outputs))

    return myNet
}