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
    Strength float64   `json:"strength"`
    From *Node         `json:"-"`
    // To *Node           `json:"-"`
    HoldingVal float64 `json:"holding"`
}

type Node struct {
    Value float64                      `json:"value"`  // should this be float32?  idk
    IncomingConnections []*Connection  `json:"connections"`  //which nodes to read from
    Position [3]int                    `json:"position"`
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
    Position [3]int   `json:"position"`
    Strength float64  `json:"strength"`
}

func (s Stimulus) String() string {
    jsonRep, _ := json.MarshalIndent(s, "", "    ")
    return string(jsonRep)
}

func (c Connection) String() string {
    jsonRep, _ := json.MarshalIndent(c, "", "    ")
    return string(jsonRep)
}

func (n *Node) Update() {
    // fmt.Println("should update")
    // figure out how to do this?
    // maybe multiply by avg. of incoming signals?
    var final float64
    for _, conn := range n.IncomingConnections {
        final = final + conn.HoldingVal
    }
    final = final / float64(len(n.IncomingConnections))
    fmt.Println(final)
}

func (net *Network) Cycle() {
    // fake concurrency
    // first, set all the connections based on their nodes
    for _, node := range net.Nodes {
        for _, conn := range node.IncomingConnections {
            conn.HoldingVal = conn.From.Value
            conn.Strength = conn.Strength * conn.From.Value
            if conn.Strength < 0.1 {
                conn.Strength = 0.1
            }
        }
    }
    // then set all the nodes based on connections
    for _, node := range net.Nodes {
        node.Update()
    }
}

func (n Node) String() string {
    jsonRep, _ := json.MarshalIndent(n, "", "    ")
    return string(jsonRep)
}

func (net *Network) Stimulate(stimuli []Stimulus) {
    for _, stim := range stimuli {
        // fmt.Println(stim)
        var applyTo *Node;
        for _, node := range net.Nodes {
            if node.Position == stim.Position {
                applyTo = node
                break
            }
        }
        applyTo.Value = stim.Strength
    }
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
                    Strength: 1,
                })
            }
        }
    }
}

func (net Network) String() string {
    jsonRep, _ := json.MarshalIndent(net, "", "    ")
    return string(jsonRep)
}

func MakeNetwork(input, perZone, output, cycles int, dimensions [3]int) *Network {
    nodes := []*Node{}
    math.Mod(5, 5)
    for i := 1; i <= dimensions[0]; i++ {
        for j := 1; j <= dimensions[1]; j++ {
            for k := 1; k <= dimensions[2]; k++ {
                nodes = append(nodes, &Node {
                    Value: 1,
                    Position: [3]int{i, j, k},
                })
            }
        }
    }
    return &Network {
        Nodes: nodes,
        CurCycle: 0,
        MaxCycles: cycles,
    }
}

func main() {
    // input, processing, output, cycles, width, depth, height, perZone
    myNet := MakeNetwork(2, 2, 1, 25, [3]int{2, 2, 2})
    myNet.Connect()
    myNet.Stimulate([]Stimulus{
        Stimulus{
            Position: [3]int{1,1,1},
            Strength: 0.5,
        },
    })
    myNet.Cycle()
    // fmt.Println(myNet)
    // fmt.Println(len(myNet.Nodes[13].IncomingConnections))
}