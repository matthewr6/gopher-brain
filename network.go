package main

import (
    "fmt"
    "math"
    "math/rand"
    "os"
    "strconv"
    "encoding/json"
)

/*
    & is "address of"
    * is "value at address"
*/

type Connection struct {
    Strength float64   `json:"strength"`
    To *Node           `json:"-"`
    HoldingVal float64 `json:"holding"`
}

type Node struct {
    Value float64                      `json:"value"`
    OutgoingConnections []*Connection  `json:"-"`  //which nodes to send to
    IncomingConnections []*Connection  `json:"-"`  //which nodes to read from
    Position [3]int                    `json:"position"`
}

type Network struct {
    Nodes []*Node `json:"nodes"`
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
    // TODO - REWORK?
    var final float64
    for _, conn := range n.IncomingConnections {
        final = final + conn.HoldingVal*conn.Strength
    }
    n.Value = final
}

func RandFloat(min, max float64) float64 {
    randFloat := rand.Float64()
    diff := max - min
    r := randFloat * diff
    return min + r
}

func (n Node) RandOutConn() *Connection {
    var ret *Connection
    sum := 0.0
    for _, conn := range n.OutgoingConnections {
        sum += conn.Strength
    }
    r := RandFloat(0.0, sum)
    for _, conn := range n.OutgoingConnections {
        r -= conn.Strength
        if r < 0 {
            return conn
        }
    }
    return ret
}

func (net *Network) Cycle() {
    // fake concurrency
    // first, set all the connections based on their nodes

    for _, node := range net.Nodes {
        outputConn := node.RandOutConn()
        outputConn.HoldingVal = node.Value
        outputConn.Strength = outputConn.Strength * node.Value // change?
        if outputConn.Strength < 0.1 || outputConn.Strength > 10 { // TWEAK THIS MAX STRENGTH!
            outputConn.Strength = 0.1
        }
        node.Value = 0 // change to node.Value * 0.25 or something?
    }

    // then set all the nodes based on connections

    for _, node := range net.Nodes {
        node.Update()
    }

    // then clear the connections
    for _, node := range net.Nodes {
        for _, conn := range node.OutgoingConnections {
            // conn.HoldingVal = conn.HoldingVal * 0.25 // what number to use
            conn.HoldingVal = 0 // use this or the one above?
        }
    }
}

func (n Node) String() string {
    jsonRep, _ := json.MarshalIndent(n, "", "    ")
    return string(jsonRep)
}

func (net *Network) Stimulate(stimuli []Stimulus) {
    for _, stim := range stimuli {
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

func (net *Network) RandomizeValues() {
    for _, node := range net.Nodes {
        node.Value = rand.Float64()
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
                newConn := &Connection{
                    To: potConNode,
                    Strength: rand.Float64() + 0.5, // do random strength - from 0.5 to 1.5?
                    // Strength: rand.Float64(), // do random strength
                }
                node.OutgoingConnections = append(node.OutgoingConnections, newConn)
                potConNode.IncomingConnections = append(potConNode.IncomingConnections, newConn)
            }
        }
    }
}

func (net Network) String() string {
    jsonRep, _ := json.MarshalIndent(net, "", "    ")
    return string(jsonRep)
}

func (net Network) DumpJSON(name string) {
    f, _ := os.Create(fmt.Sprintf("./net_%v.json", name))
    f.WriteString(net.String())
    f.Close()
}

func MakeNetwork(dimensions [3]int, blank bool) *Network {
    nodes := []*Node{}
    math.Mod(5, 5)
    for i := 1; i <= dimensions[0]; i++ {
        for j := 1; j <= dimensions[1]; j++ {
            for k := 1; k <= dimensions[2]; k++ {
                var newValue float64
                if !blank {
                    newValue = rand.Float64()
                }
                nodes = append(nodes, &Node {
                    Value: newValue,
                    Position: [3]int{i, j, k},
                })
            }
        }
    }
    return &Network {
        Nodes: nodes,
    }
}

func (net *Network) GenerateAnim(frames int) {
    for frame := 0; frame < frames; frame++ {
        net.DumpJSON(strconv.Itoa(frame))
        net.Cycle()
    }
}