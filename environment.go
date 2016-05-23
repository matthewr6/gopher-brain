package main

import (
    "math/rand"
    "encoding/json"
    // "fmt"
)

// TODO - some sort of equation fitted on each to determine the response
// TODO - some sort of responder struct - should it be a many-to-many relationship?
//      - if so - one receiver can influence many whatevers, and one whatever can be influenced by multiple receivers
// dangit gonna have to add this to savestate/loadstate

// sensors feed data to nodes
type Sensor struct {
    Nodes []*Node       `json:"nodes"`
    Excitatory bool     `json:"excitatory"` // todo this probably isn't used
    Trigger string      `json:"trigger"`
    Stimulated bool     `json:"-"`
    Name string         `json:"name"`
}

type Output struct {
    Nodes []*Node   `json:"nodes"`
    Name string     `json:"name"`
    Value float64   `json:"value"` // do we want single or double precision float?  (32/64)
}

func (s Sensor) String() string {
    jsonRep, _ := json.MarshalIndent(s, "", "    ")
    return string(jsonRep)
}

func (o Output) String() string {
    jsonRep, _ := json.MarshalIndent(o, "", "    ")
    return string(jsonRep)
}

func (output *Output) Update() {
    var sum float64
    for _, node := range output.Nodes {
        if node.OutgoingConnection.Excitatory {
            sum += float64(node.Value) * node.OutgoingConnection.Strength
        } else {
            sum -= float64(node.Value) * node.OutgoingConnection.Strength
        }
    }
    output.Value = sum
}

func (sensor *Sensor) Update() {
    // for now let's just continuously stimulate every node
    for _, node := range sensor.Nodes {
        if sensor.Stimulated {
            node.Value = 1
        } else {
            node.Value = 0
        }
    }
}

// todo - there's probably an easier way to do the plane stuff now

// do I even need the plane stuff
// seems bloated
// todo reorder these args
// also it's SO LONG AND MESSY :L
func (net *Network) CreateSensor(name string, r int, count int, plane string, center [3]int, excitatory bool, trigger string) *Sensor {
    // radius is basically density...
    sensor := &Sensor{
        // Radius: r,
        // NodeCount: count,
        Nodes: []*Node{},
        Excitatory: excitatory,
        Trigger: trigger,
        Stimulated: false,
        Name: name,
        // Center: center,
    }
    // if kb != nil {
    //     kb.Bind(func() {
    //         sensor.Stimulated = !sensor.Stimulated
    //     }, trigger)
    // }
    // todo - determine correct coefficient
    stDev := float64(r)
    // plane is which dimension should stay the same - name the variable in a better way?
    if (plane != "") {
        if (plane == "x" || plane == "y" || plane == "z") {
            // todo - also this coefficient
            stDev = float64(r * 2)
        }
        if (plane == "x") {
            potX := center[0]
            for len(sensor.Nodes) < count {
                potY := int(rand.NormFloat64() * stDev) + center[1]
                potZ := int(rand.NormFloat64() * stDev) + center[2]
                if potY > 0 && potZ > 0 {
                    potNode := FindNode([3]int{potX, potY, potZ}, net.Nodes)
                    if !NodeExistsIn(potNode, sensor.Nodes) {
                        sensor.Nodes = append(sensor.Nodes, potNode)
                    }
                }
            }
        }
        if (plane == "y") {
            potY := center[1]
            for len(sensor.Nodes) < count {
                potX := int(rand.NormFloat64() * stDev) + center[0]
                potZ := int(rand.NormFloat64() * stDev) + center[2]
                if potX > 0 && potZ > 0 {
                    potNode := FindNode([3]int{potX, potY, potZ}, net.Nodes)
                    if !NodeExistsIn(potNode, sensor.Nodes) {
                        sensor.Nodes = append(sensor.Nodes, potNode)
                    }
                }
            }
        }
        if (plane == "z") {
            potZ := center[2]
            for len(sensor.Nodes) < count {
                potX := int(rand.NormFloat64() * stDev) + center[0]
                potY := int(rand.NormFloat64() * stDev) + center[1]
                if potX > 0 && potY > 0 {
                    potNode := FindNode([3]int{potX, potY, potZ}, net.Nodes)
                    if !NodeExistsIn(potNode, sensor.Nodes) {
                        sensor.Nodes = append(sensor.Nodes, potNode)
                    }
                }
            }
        }
    } else {
        for len(sensor.Nodes) < count {
            potX := int(rand.NormFloat64() * stDev) + center[0]
            potY := int(rand.NormFloat64() * stDev) + center[1]
            potZ := int(rand.NormFloat64() * stDev) + center[2]
            if potX >= 0 && potY >= 0 && potZ >= 0 && potX < net.Dimensions[0] && potY < net.Dimensions[1] && potZ < net.Dimensions[2] {
                potNode := FindNode([3]int{potX, potY, potZ}, net.Nodes)
                if !NodeExistsIn(potNode, sensor.Nodes) {
                    sensor.Nodes = append(sensor.Nodes, potNode)
                }
            }
        }
    }
    net.Sensors = append(net.Sensors, sensor)
    return sensor
}

func (net *Network) CreateOutput(name string, r int, count int, plane string, center [3]int) *Output {
    // radius is basically density...
    output := &Output{
        Nodes: []*Node{},
        Name: name,
        // Center: center,
    }
    // todo - determine correct coefficient
    stDev := float64(r)
    // plane is which dimension should stay the same - name the variable in a better way?
    if (plane != "") {
        if (plane == "x" || plane == "y" || plane == "z") {
            // todo - also this coefficient
            stDev = float64(r * 2)
        }
        if (plane == "x") {
            potX := center[0]
            for len(output.Nodes) < count {
                potY := int(rand.NormFloat64() * stDev) + center[1]
                potZ := int(rand.NormFloat64() * stDev) + center[2]
                if potY > 0 && potZ > 0 {
                    potNode := FindNode([3]int{potX, potY, potZ}, net.Nodes)
                    if !NodeExistsIn(potNode, output.Nodes) {
                        output.Nodes = append(output.Nodes, potNode)
                    }
                }
            }
        }
        if (plane == "y") {
            potY := center[1]
            for len(output.Nodes) < count {
                potX := int(rand.NormFloat64() * stDev) + center[0]
                potZ := int(rand.NormFloat64() * stDev) + center[2]
                if potX > 0 && potZ > 0 {
                    potNode := FindNode([3]int{potX, potY, potZ}, net.Nodes)
                    if !NodeExistsIn(potNode, output.Nodes) {
                        output.Nodes = append(output.Nodes, potNode)
                    }
                }
            }
        }
        if (plane == "z") {
            potZ := center[2]
            for len(output.Nodes) < count {
                potX := int(rand.NormFloat64() * stDev) + center[0]
                potY := int(rand.NormFloat64() * stDev) + center[1]
                if potX > 0 && potY > 0 {
                    potNode := FindNode([3]int{potX, potY, potZ}, net.Nodes)
                    if !NodeExistsIn(potNode, output.Nodes) {
                        output.Nodes = append(output.Nodes, potNode)
                    }
                }
            }
        }
    } else {
        for len(output.Nodes) < count {
            potX := int(rand.NormFloat64() * stDev) + center[0]
            potY := int(rand.NormFloat64() * stDev) + center[1]
            potZ := int(rand.NormFloat64() * stDev) + center[2]
            if potX >= 0 && potY >= 0 && potZ >= 0 && potX < net.Dimensions[0] && potY < net.Dimensions[1] && potZ < net.Dimensions[2] {
                potNode := FindNode([3]int{potX, potY, potZ}, net.Nodes)
                if !NodeExistsIn(potNode, output.Nodes) {
                    output.Nodes = append(output.Nodes, potNode)
                }
            }
        }
    }
    net.Outputs = append(net.Outputs, output)
    return output
}