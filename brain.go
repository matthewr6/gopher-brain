package brain

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

func Brain(NETWORK_SIZE [3]int, CONSTRUCTORS []SensorConstructor) {
    reader := bufio.NewReader(os.Stdin)
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
    var mode string
    tracker := make(map[string]bool)
    fmt.Printf("Currently has %v sensors and %v outputs.\n", len(myNet.Sensors), len(myNet.Outputs))
    fmt.Println("Sensor names:")
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
        // let's pretend the front x/z plane (y = 1) is "front" with left being x = 25
        // maybe you should only create sensors, and specify # of corresponding outputs - and then the createSensor generates the outputs automatically
        for _, constructor := range CONSTRUCTORS {
            myNet.CreateSensor(constructor.Name, constructor.R, constructor.Count, constructor.Plane, constructor.Center, constructor.OutputCount, constructor.InputFunc)
        }
    }
    myNet.PruneUnusedSensors()
    fmt.Printf("Now has %v sensors and %v outputs.\n", len(myNet.Sensors), len(myNet.Outputs))

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
    
    myNet.GenerateAnim(frames)

    fmt.Print("\nSave state?  Enter a name if you wish to save the state:  ")
    fileName, _ = reader.ReadString('\n')
    fileName = strings.TrimSpace(fileName)
    if fileName != "" {
        myNet.SaveState(fileName)
    }
}