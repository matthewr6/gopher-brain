package main

import (
    "math/rand"
    "encoding/json"
)

// TODO - some sort of equation fitted on each to determine the response
// TODO - some sort of responder struct - should it be a many-to-many relationship?
//      - if so - one receiver can influence many whatevers, and one whatever can be influenced by multiple receivers
// dangit gonna have to add this to savestate/loadstate

// sensors feed data to nodes
// todo what attrs do I need
type Sensor struct {
    // Radius int          `json:"radius"`
    // NodeCount int       `json:"nodeCount"`
    Nodes []*Node       `json:"nodes"`
    Excitatory bool     `json:"excitatory"`
    Trigger string      `json:"trigger"`
    // Center [3]int       `json:"center"`
    // FlatPlane string    `json:"plane"` // make it a flat plane
}

func (s Sensor) String() string {
    jsonRep, _ := json.MarshalIndent(s, "", "    ")
    return string(jsonRep)
}

func (sensor *Sensor) Update() {
    // for now let's just continuously stimulate every node
    // maybe try randomly lighting up the node, a 50/50 chance?
    for _, node := range sensor.Nodes {
        if (sensor.Excitatory) {
            node.Value = 1
        } else {
            node.Value = 0
        }
    }
}

// do I even need the plane stuff
// seems bloated
// todo reorder these args
func (net *Network) CreateSensor(r int, count int, plane string, center [3]int, excitatory bool, trigger string) *Sensor {
    // radius is basically density...
    sensor := &Sensor{
        // Radius: r,
        // NodeCount: count,
        Nodes: []*Node{},
        Excitatory: excitatory,
        Trigger: trigger,
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
            if potX > 0 && potY > 0 && potZ > 0 {
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