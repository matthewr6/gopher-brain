package main

import (
    "os"
    "fmt"
    "time"
    "bufio"
    "strings"
    "math/rand"

    "github.com/jteeuwen/keyboard/termbox"
    term "github.com/nsf/termbox-go"
)

/*
    & is "address of"
    * is "value at address"
*/

var running = true

func main() {
    reader := bufio.NewReader(os.Stdin)
    fileName := Prompt("Enter state name to load state, or leave blank to create a new network:  ", reader)

    rand.Seed(time.Now().UTC().UnixNano())

    var myNet *Network
    _, err := os.Stat(fmt.Sprintf("./state/%v_state.json", fileName))
    if fileName == "" || err != nil {
        NETWORK_SIZE := [3]int{3, 3, 3}
        myNet = MakeNetwork(NETWORK_SIZE, false)
        myNet.Connect()
        myNet.Mirror()
    } else {
        myNet = LoadState(fileName)
    }
    var mode string
    mode = Prompt("Custom or simple?  Defaults to simple.  [custom/simple]  ", reader)
    if mode == "custom" {
        fmt.Println("WARNING!  Sensors and outputs will not save properly in this mode!")
        myNet.ClearIO()
        // let's pretend the front x/z plane (y = 1) is "front" with left being x = 25
        myNet.CreateSensor("left (a)", 1, 50, "y", [3]int{24, 0, 12}, true, "a", func(nodes []*Node, stimulated bool) {

        })
        myNet.CreateSensor("right (d)", 1, 50, "y", [3]int{0, 0, 12}, true, "d", func(nodes []*Node, stimulated bool) {

        })

        // myNet.CreateOutput("front")
        // myNet.CreateOutput("left")
        // myNet.CreateOutput("right")
        // myNet.CreateOutput("back")
    } else {
        var choice string
        fmt.Printf("\nNetwork has %v sensor(s):\n", len(myNet.Sensors))
        for _, sensor := range myNet.Sensors {
            fmt.Printf("    %v\n", sensor.Name)
        }
        choice = Prompt("\nAdd sensor? [y/n]  ", reader)
        for choice == "y" {
            sensorName := Prompt("    Name:  ", reader)
            trigger := Prompt("    Trigger [single key]:  ", reader) // should validate to be one key
            plane := Prompt("    Plane [x/y/z/blank]:  ", reader)
            if plane != "x" && plane != "y" && plane != "z" {
                plane = ""
            }
            // todo - validate for negatives
            centerArr := []int{}
            for len(centerArr) != 3 {
                center := Prompt("    Center [format x,y,z]:  ", reader)
                centerArr = StrsToInts(strings.Split(center, ","))
            }
            // todo find numbers and stuff
            myNet.CreateSensor(sensorName, 1, 50, plane, [3]int{centerArr[0], centerArr[1], centerArr[2]}, true, trigger, func(nodes []*Node, stimulated bool) {
                // for simplicity - just continuously stimulate every node
                for _, node := range nodes {
                    if stimulated {
                        node.Value = 1
                    }
                    // let's try removing this for now, see what happens...
                    // else {
                    //     node.Value = 0
                    // }
                }
            })
            choice = Prompt("\nAdd another sensor? [y/n]  ", reader)
        }
        choice = Prompt("\nEnter a sensor name to remove a sensor:  ", reader)
        for choice != "" {
            myNet.RemoveSensor(choice)
            choice = Prompt("Enter another sensor name to remove:  ", reader)
        }
        fmt.Printf("\nNetwork has %v output(s).\n", len(myNet.Outputs))
        for _, output := range myNet.Outputs {
            fmt.Printf("    %v\n", output.Name)
        }
        choice = Prompt("\n    Add output? [y/n]  ", reader)
        for choice == "y" {
            outputName := Prompt("    Name:  ", reader)
            plane := Prompt("    Plane [x/y/z/blank]:  ", reader)
            if plane != "x" && plane != "y" && plane != "z" {
                plane = ""
            }
            // todo - validate for negatives
            centerArr := []int{}
            for len(centerArr) != 3 {
                center := Prompt("    Center [format x,y,z]:  ", reader)
                centerArr = StrsToInts(strings.Split(center, ","))
            }
            // todo get numbers
            myNet.CreateOutput(outputName, 1, 50, plane, [3]int{centerArr[0], centerArr[1], centerArr[2]}, func(nodes map[*Node]*ConnInfo) float64 {
                var sum float64
                for node, connInfo := range nodes {
                    if connInfo.Excitatory {
                        sum += float64(node.Value) * connInfo.Strength
                    } else {
                        sum -= float64(node.Value) * connInfo.Strength
                    }
                }
                return sum
            })
            choice = Prompt("Add another output? [y/n]  ", reader)
        }
        choice = Prompt("\nEnter an output name to remove an output:  ", reader)
        for choice != "" {
            myNet.RemoveOutput(choice)
            choice = Prompt("Enter another output name to remove:  ", reader)
        }
    }

    // this is the keyboard sensing stuff
    term.Init()
    term.SetCursor(0, 0)
    kb := termbox.New()
    kb.Bind(func() {
        running = false
    }, "space")
    go KeyboardPoll(kb)
    myNet.BindKeyboard(kb)

    myNet.AnimateUntilDone(100)

    term.Close()

    fmt.Print("\nSave state?  Enter a name if you wish to save the state:  ")
    fileName, _ = reader.ReadString('\n')
    fileName = strings.TrimSpace(fileName)
    if fileName != "" {
        myNet.SaveState(fileName)
    }
}