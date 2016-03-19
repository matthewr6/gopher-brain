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
    Value float64               `json:"value"`  // should this be float32?  idk
    IncomingConnections []*Connection  `json:"connections"`  //which nodes to read from
    // Level int                   `json:"level"`  // for 
    Position [3]int             `json:"position"`
}

type Network struct {
    // InputNodes []*Node   `json:"inputNodes"`
    Nodes []*Node        `json:"nodes"`
    // OutputNodes []*Node  `json:"outputNodes"`
    MaxLevel int         `json:"maxLevel"`
    MaxCycles int        `json:"maxCycles"`
    CurCycle int         `json:"curCycle"`
}

type Stimulus struct {
    Position [3]int
    Strength float64
}

func (n *Node) Update() {
    // figure out how to do this
    // synapses should strengthen when used and weaken when not - maybe multiply by .9 and 1.1 or something?
    // go by levels - so it can kinda radiate outwards
}

func (n Node) String() string {
    jsonRep, _ := json.MarshalIndent(n, "", "    ")
    return string(jsonRep)
}

func (net *Network) Stimulate([]Stimulus) {

}

func ThreeDimDist(p1, p2 [3]int) float64 {
    ans := (p1[0]-p2[0])*(p1[0]-p2[0]) + (p1[1]-p2[1])*(p1[1]-p2[1]) + (p1[2]-p2[2])*(p1[2]-p2[2])
    return math.Sqrt(float64(ans))
}
func (net *Network) Connect() {
    for _, node := range net.Nodes {
        // get the closest nodes and connect
        for _, potConNode := range net.Nodes {
            if ThreeDimDist(node.Position, potConNode.Position) < 1.75 && node != potConNode {
                node.IncomingConnections = append(node.IncomingConnections, &Connection {
                    From: potConNode,
                    Strength: 0.5,
                })
            }
        }
    }
    // // first connect all "body" level 1 nodes to "input"
    // for _, bodyNode := range net.Nodes {
    //     if bodyNode.Level == 1 {
    //         for _, inputNode := range net.InputNodes {
    //             bodyNode.Connections = append(bodyNode.Connections, &Connection{
    //                 Strength: 0.5,
    //                 From: inputNode,
    //             })
    //         }
    //     }
    // }
    // // then connect all "body" nodes to ones close to them
    // for i := 1; i <= net.MaxLevel; i++ {
    //     for _, bodyNode := range net.Nodes {
    //         if bodyNode.Level == i {
    //             for _, otherNode := range net.Nodes {
    //                 if otherNode.Level == i && otherNode != bodyNode {
    //                     bodyNode.Connections = append(bodyNode.Connections, &Connection{
    //                         Strength: 0.5,
    //                         From: otherNode,
    //                     })
    //                 }
    //             }
    //         }
    //     }
    // }
    // then connect all "body" nodes to level below them

    // then connect all "output" nodes to "body"

}

func (net *Network) Cycle() {
    for _, node := range net.Nodes {
        fmt.Println(node)
    }
    net.CurCycle++
}

func (net Network) String() string {
    jsonRep, _ := json.MarshalIndent(net, "", "    ")
    return string(jsonRep)
}

func MakeNetwork(input, perZone, output, cycles int, dimensions [3]int) *Network {
    nodes := []*Node{}
    // inputNodes := []*Node{}
    // outputNodes := []*Node{}
    // for i := 0; i < input; i++ {
    //     inputNodes = append(inputNodes, &Node{
    //         Value: 0,
    //         // Level: 0,
    //     })
    // }
    math.Mod(5, 5)
    for i := 1; i <= dimensions[0]; i++ {
        for j := 1; j <= dimensions[1]; j++ {
            for k := 1; k <= dimensions[2]; k++ {
                nodes = append(nodes, &Node {
                    Value: 0,
                    Position: [3]int{i, j, k},
                })
            }
        }
    }

    // for i := 0; i < output; i++ {
    //     outputNodes = append(outputNodes, &Node{
    //         Value: 0,
    //         // Level: curLevel + 1, // change this based on above comment w/leveling
    //     })
    // }
    return &Network {
        Nodes: nodes,
        // InputNodes: inputNodes,
        // OutputNodes: outputNodes,
        // MaxLevel: curLevel,
        CurCycle: 0,
        MaxCycles: cycles,
    }
}

func main() {
    // input, processing, output, cycles, width, depth, height, perZone
    myNet := MakeNetwork(2, 2, 1, 25, [3]int{3, 3, 3})
    myNet.Connect()
    // myNet.Cycle()
    // fmt.Println(myNet)
    // fmt.Println(len(myNet.Nodes[13].IncomingConnections))
}