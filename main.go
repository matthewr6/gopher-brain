package main

import (
    "fmt"
    "math"
    "encoding/json"

    // "reflect"
)

/*
    & is "address of"
    * is "value at address"
*/

type Node struct {
    Value float64        `json:"value"`  // should this be float32?  idk
    Connections []*Node  `json:"connections"`  //which nodes to read from
    Level int            `json:"level"`  // for 
}

type Network struct {
    InputNodes []*Node   `json:"inputNodes"`
    Nodes []*Node        `json:"nodes"`
    OutputNodes []*Node  `json:"outputNodes"`
    MaxCycles int        `json:"maxCycles"`
    CurCycle int         `json:"curCycle"`
}

func (n *Node) update() {
    // figure out how to do this
    // two approaches - all connect all, or a random X number of connections
    // random would probably be a little less predictable, so probably better to do random
}

func (n Node) String() string {
    jsonRep, _ := json.MarshalIndent(n, "", "    ")
    return string(jsonRep)
}

func (net *Network) connect() {
    // do things?
}

func (net *Network) cycle() {
    for _, node := range net.Nodes {
        fmt.Println(node)
    }
    net.CurCycle++
}

func (net Network) String() string {
    jsonRep, _ := json.MarshalIndent(net, "", "    ")
    return string(jsonRep)
}

func MakeNetwork(input, processing, output, cycles, perLevel int) *Network {
    nodes := []*Node{}
    inputNodes := []*Node{}
    outputNodes := []*Node{}
    for i := 0; i < input; i++ {
        inputNodes = append(nodes, &Node{
            Value: 0,
            Level: 0,
        })
    }
    curLevel := 1
    for i := 1; i <= processing; i++ {
        nodes = append(nodes, &Node{
            Value: 0,
            Connections: []*Node{},
            Level: curLevel,
        })
        if math.Mod(float64(i), float64(perLevel)) == 0 {
            curLevel++
        }
    }
    for i := 0; i < output; i++ {
        outputNodes = append(outputNodes, &Node{
            Value: 0,
            Level: curLevel + 1, // change this based on above comment w/leveling
        })
    }
    return &Network {
        Nodes: nodes,
        InputNodes: inputNodes,
        OutputNodes: outputNodes,
        CurCycle: 0,
        MaxCycles: cycles,
    }
}

func main() {
    // input, processing, output, cycles, perLevel
    myNet := MakeNetwork(1, 5, 1, 25, 2)
    myNet.connect()
    myNet.cycle()
    fmt.Println(myNet)
}