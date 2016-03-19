package main

import (
    "fmt"
    "encoding/json"
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

func (net Network) print() {
    jsonRep, _ := json.MarshalIndent(net, "", "    ")
    fmt.Println(string(jsonRep))
}

func MakeNetwork(input, processing, output, cycles, perLevel int) *Network {
    nodes := []*Node{}
    inputNodes := []*Node{}
    outputNodes := []*Node{}
    for i := 0; i < processing; i++ {
        nodes = append(nodes, &Node{
            Value: 0,
            Level: 0,
        })
    }
    baseLevel := 1
    for i := 0; i < input; i++ {
        inputNodes = append(inputNodes, &Node{
            Value: 1,
            Connections: []*Node{},
            Level: baseLevel,
        })
        // if i mod perlevel is 0, increase baselevel?
    }
    for i := 0; i < output; i++ {
        outputNodes = append(outputNodes, &Node{
            Value: 0,
            Level: input + 1, // change this based on above comment w/leveling
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
    myNet := MakeNetwork(5,2,6, 25, 5)
    myNet.connect()
    myNet.cycle()
    // myNet.print()
}