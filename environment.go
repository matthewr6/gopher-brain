package main

import (
    "math/rand"
    "encoding/json"
)

// sensors feed data to nodes
type Sensor struct {
    Nodes []*Node          `json:"nodes"`
    Excitatory bool        `json:"excitatory"` // todo this probably isn't used
    Trigger string         `json:"trigger"`
    Name string            `json:"name"`
    In func([]*Node)       `json:"-"`
}

// dang gonna have to do the same saving trick stuff as the Connection type
type Output struct {
    Nodes map[*Node]*ConnInfo               `json:"nodes"`
    Name string                             `json:"name"`
    Value float64                           `json:"value"`
    Out func(map[*Node]*ConnInfo) float64   `json:"-"`
}

func (s Sensor) String() string {
    jsonRep, _ := json.MarshalIndent(s, "", "    ")
    return string(jsonRep)
}

func (o Output) String() string {
    jsonRep, _ := json.MarshalIndent(o, "", "    ")
    return string(jsonRep)
}

func (net *Network) ClearIO() {
    net.RemoveAllSensors()
    net.RemoveAllOutputs()
}

func (net *Network) RemoveAllSensors() {
    net.Sensors = []*Sensor{}
}

func (net *Network) RemoveAllOutputs() {
    net.Outputs = []*Output{}
}

func (net *Network) RemoveSensor(name string) {
    index := len(net.Sensors)
    for i, sensor := range net.Sensors {
        if sensor.Name == name {
            index = i
            break
        }
    }
    if index != len(net.Sensors) {
        net.Sensors = append(net.Sensors[:index], net.Sensors[index+1:]...)
    }
}

func (net *Network) RemoveOutput(name string) {
    index := len(net.Outputs)
    for i, output := range net.Outputs {
        if output.Name == name {
            index = i
            break
        }
    }
    if index != len(net.Outputs) {
        net.Outputs = append(net.Outputs[:index], net.Outputs[index+1:]...)
    }
}

// todo - there's probably an easier way to do the plane stuff now

// do I even need the plane stuff
// seems bloated
// todo reorder these args
// also it's SO LONG AND MESSY :L
func (net *Network) CreateSensor(name string, r int, count int, plane string, center [3]int, excitatory bool, trigger string, inputFunc func([]*Node)) *Sensor {
    // radius is basically density...
    sensor := &Sensor{
        Nodes: []*Node{},
        Excitatory: excitatory,
        Trigger: trigger,
        Name: name,
        In: inputFunc,
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
            for len(sensor.Nodes) < count {
                potY := int(rand.NormFloat64() * stDev) + center[1]
                potZ := int(rand.NormFloat64() * stDev) + center[2]
                if potY > 0 && potZ > 0 {
                    potNode := net.FindNode([3]int{potX, potY, potZ})
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
                    potNode := net.FindNode([3]int{potX, potY, potZ})
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
                    potNode := net.FindNode([3]int{potX, potY, potZ})
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
                potNode := net.FindNode([3]int{potX, potY, potZ})
                if !NodeExistsIn(potNode, sensor.Nodes) {
                    sensor.Nodes = append(sensor.Nodes, potNode)
                }
            }
        }
    }
    net.Sensors = append(net.Sensors, sensor)
    return sensor
}

func (net *Network) CreateOutput(name string, r int, count int, plane string, center [3]int, outputFunc func(map[*Node]*ConnInfo) float64) *Output {
    // radius is basically density...
    output := &Output{
        Name: name,
        Out: outputFunc,
    }
    // todo - determine correct coefficient
    stDev := float64(r)

    // set up nodes
    nodes := []*Node{}
    // plane is which dimension should stay the same - name the variable in a better way?
    if (plane != "") {
        if (plane == "x" || plane == "y" || plane == "z") {
            // todo - also this coefficient
            stDev = float64(r * 2)
        }
        if (plane == "x") {
            potX := center[0]
            for len(nodes) < count {
                potY := int(rand.NormFloat64() * stDev) + center[1]
                potZ := int(rand.NormFloat64() * stDev) + center[2]
                if potY > 0 && potZ > 0 {
                    potNode := net.FindNode([3]int{potX, potY, potZ})
                    if !NodeExistsIn(potNode, nodes) {
                        nodes = append(nodes, potNode)
                    }
                }
            }
        }
        if (plane == "y") {
            potY := center[1]
            for len(nodes) < count {
                potX := int(rand.NormFloat64() * stDev) + center[0]
                potZ := int(rand.NormFloat64() * stDev) + center[2]
                if potX > 0 && potZ > 0 {
                    potNode := net.FindNode([3]int{potX, potY, potZ})
                    if !NodeExistsIn(potNode, nodes) {
                        nodes = append(nodes, potNode)
                    }
                }
            }
        }
        if (plane == "z") {
            potZ := center[2]
            for len(nodes) < count {
                potX := int(rand.NormFloat64() * stDev) + center[0]
                potY := int(rand.NormFloat64() * stDev) + center[1]
                if potX > 0 && potY > 0 {
                    potNode := net.FindNode([3]int{potX, potY, potZ})
                    if !NodeExistsIn(potNode, nodes) {
                        nodes = append(nodes, potNode)
                    }
                }
            }
        }
    } else {
        for len(nodes) < count {
            potX := int(rand.NormFloat64() * stDev) + center[0]
            potY := int(rand.NormFloat64() * stDev) + center[1]
            potZ := int(rand.NormFloat64() * stDev) + center[2]
            if potX >= 0 && potY >= 0 && potZ >= 0 && potX < net.Dimensions[0] && potY < net.Dimensions[1] && potZ < net.Dimensions[2] {
                potNode := net.FindNode([3]int{potX, potY, potZ})
                if !NodeExistsIn(potNode, nodes) {
                    nodes = append(nodes, potNode)
                }
            }
        }
    }

    // iterate through nodes
    nodeMapping := make(map[*Node]*ConnInfo)
    var excitatory bool
    for _, node := range nodes {
        if rand.Intn(2) != 0 {
            excitatory = true
        }
        nodeMapping[node] = &ConnInfo{
            Strength: RandFloat(0.75, 1.75),
            Excitatory: excitatory,
        }
    }
    output.Nodes = nodeMapping
    net.Outputs = append(net.Outputs, output)
    return output
}