package main

import (
    "fmt"
    "math/rand"
    "encoding/json"
)

// sensors feed data to nodes
/*
    sensor = input
    one sensor feeds to a set of nodes
    one output is only influenced by its list of connected nodes
    one sensor can be influenced by a SET of outputs
*/
type Sensor struct {
    Nodes []*Node                   `json:"nodes"`
    Influences []*Output            `json:"influences"`
    Name string                     `json:"name"`
    In func([]*Node, []*Output)     `json:"-"`
}

// do I want to save sensors/outputs?
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
    net.Sensors = make(map[string]*Sensor)
}

func (net *Network) RemoveAllOutputs() {
    net.Outputs = make(map[string]*Output)
}

func (net *Network) RemoveSensor(name string) {
    delete(net.Sensors, name)
}

func (net *Network) RemoveOutput(name string) {
    delete(net.Outputs, name)
}

// todo - there's probably an easier way to do the plane stuff now

// do I even need the plane stuff
// seems bloated
// todo reorder these args
// also it's SO LONG AND MESSY :L
func (net *Network) CreateSensor(name string, r int, count int, plane string, center [3]int, outputCount int, inputFunc func([]*Node, []*Output)) [2]*Sensor {
    secondCenter := center
    secondCenter[0] = (net.Dimensions[0]*2) - center[0]
    // set random output centers here
    outputCenters := [][3]int{}
    for i := 0; i < outputCount; i++ {
        outputCenters = append(outputCenters, [3]int{
            rand.Intn(net.Dimensions[0]),
            rand.Intn(net.Dimensions[1]),
            rand.Intn(net.Dimensions[2]),
        })
    }
    a := net.CreateIndividualSensor(fmt.Sprintf("%v-one", name), r, count, plane, center, true, outputCenters, inputFunc)
    b := net.CreateIndividualSensor(fmt.Sprintf("%v-two", name), r, count, plane, secondCenter, false, outputCenters, inputFunc)
    return [2]*Sensor{a, b}
}

func (net *Network) MakeOutputs(sensorName string, outputCenters [][3]int, r int, count int, otherSide bool) []*Output {
    outputs := []*Output{}
    // for i := 0; i < outputCenters; i++ {
    // }
    for idx, center := range outputCenters {
        outputCenter := center
        if otherSide {
            outputCenter[0] = (net.Dimensions[0]*2) - outputCenter[0]
        }
        newOutput := net.CreateIndividualOutput(fmt.Sprintf("%v-%v", sensorName, idx), r, count, "", outputCenter, func(nodes map[*Node]*ConnInfo) float64 {
            // todo - I really should have a specified ConnInfo for the output
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
        outputs = append(outputs, newOutput)
    }
    return outputs
}

func (net *Network) CreateIndividualSensor(name string, r int, count int, plane string, center [3]int, otherSide bool, outputCenters [][3]int, inputFunc func([]*Node, []*Output)) *Sensor {
    outputs := net.MakeOutputs(name,  outputCenters, r, count, otherSide)
    // radius is basically density...
    sensor := &Sensor{
        Nodes: []*Node{},
        Name: name,
        In: inputFunc,
        Influences: outputs,
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
                if potY > 0 && potZ > 0 && potY < net.Dimensions[1] && potZ < net.Dimensions[2] {
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
                if potX > 0 && potZ > 0 && potX < net.Dimensions[0]*2 && potZ < net.Dimensions[2] {
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
                if potX > 0 && potY > 0 && potX < net.Dimensions[0]*2 && potY < net.Dimensions[1] {
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
            if potX >= 0 && potY >= 0 && potZ >= 0 && potX < (net.Dimensions[0]*2) && potY < net.Dimensions[1] && potZ < net.Dimensions[2] {
                potNode := net.FindNode([3]int{potX, potY, potZ})
                if !NodeExistsIn(potNode, sensor.Nodes) {
                    sensor.Nodes = append(sensor.Nodes, potNode)
                }
            }
        }
    }
    net.Sensors[name] = sensor
    return sensor
}

func (net *Network) CreateIndividualOutput(name string, r int, count int, plane string, center [3]int, outputFunc func(map[*Node]*ConnInfo) float64) *Output {
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
            if potX >= 0 && potY >= 0 && potZ >= 0 && potX < (net.Dimensions[0]*2) && potY < net.Dimensions[1] && potZ < net.Dimensions[2] {
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
    net.Outputs[name] = output
    return output
}