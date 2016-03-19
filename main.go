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

type Connection struct {
    Strength float64 `json:"strength"`
    From *Node       `json:"-"`
    // To *Node         `json:"-"`
}

type Node struct {
    Value float64             `json:"value"`  // should this be float32?  idk
    Connections []*Connection  `json:"connections"`  //which nodes to read from
    Level int                  `json:"level"`  // for 
}

type Network struct {
    InputNodes []*Node   `json:"inputNodes"`
    Nodes []*Node        `json:"nodes"`
    OutputNodes []*Node  `json:"outputNodes"`
    MaxLevel int         `json:"maxLevel"`
    MaxCycles int        `json:"maxCycles"`
    CurCycle int         `json:"curCycle"`
}

func (n *Node) update() {
    // figure out how to do this
    // synapses should strengthen when used and weaken when not - maybe multiply by .9 and 1.1 or something?
    // go by levels - so it can kinda radiate outwards
}

func (n Node) String() string {
    jsonRep, _ := json.MarshalIndent(n, "", "    ")
    return string(jsonRep)
}

func (net *Network) connect() {
    // do things?
    // first connect all "body" nodes to "input"
    for _, bodyNode := range net.Nodes {
        bodyNode.Connections = append(bodyNode.Connections, &Connection{
            Strength: 0.5,
        })
    }
    // then connect all "body" nodes to the ones on the same level

    // then connect all "output" nodes to "body"

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

func MakeNetwork(input, processing, output, cyclesPerLevel, perLevel int) *Network {
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
        MaxLevel: curLevel,
        CurCycle: 0,
        MaxCycles: cyclesPerLevel,
    }
}

func main() {
    // input, processing, output, cyclesPerLevel, perLevel
    myNet := MakeNetwork(1, 5, 1, 25, 2)
    myNet.connect()
    // myNet.cycle()
    fmt.Println(myNet)
}