package main

import (
    "math/rand"
    "encoding/json"
)

// TODO - decide name (i.e. should it be "receiver" or not)
// TODO - some sort of equation fitted on each to determine the response
// TODO - some sort of responder struct - should it be a many-to-many relationship?
type Receiver struct {
    Radius int    `json:"radius"`
    NodeCount int `json:"nodeCount"`
    Nodes []*Node `json:"nodes"`
    Center [3]int `json:"center"`
}

func (s Receiver) String() string {
    jsonRep, _ := json.MarshalIndent(s, "", "    ")
    return string(jsonRep)
}

func NodeExistsIn(node *Node, nodes []*Node) bool {
    for _, potNode := range nodes {
        if (node == potNode) {
            return true
        }
    }
    return false
}

func (net *Network) CreateReceiver(r int, count int, center [3]int) *Receiver {
    receiver := &Receiver{
        Radius: r,
        NodeCount: count,
        Nodes: []*Node{},
        Center: center,
    }
    // todo - determine correct coefficient
    stDev := float64(r * 2)
    for len(receiver.Nodes) < count {
        potX := int(rand.NormFloat64() * stDev) + center[0]
        potY := int(rand.NormFloat64() * stDev) + center[1]
        potZ := int(rand.NormFloat64() * stDev) + center[2]
        if potX > 0 && potY > 0 && potZ > 0 {
            potNode := FindNode([3]int{potX, potY, potZ}, net.Nodes)
            if !NodeExistsIn(potNode, receiver.Nodes) {
                receiver.Nodes = append(receiver.Nodes, potNode)
            }
        }
    }
    net.Receivers = append(net.Receivers, receiver)
    return receiver
}