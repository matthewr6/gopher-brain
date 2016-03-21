package main

import (
    "fmt"
    "math"
    "math/rand"
    "os"
    "time"
    "strconv"
    "encoding/json"
)

/*
    & is "address of"
    * is "value at address"
*/

type Connection struct {
    Strength float64   `json:"strength"`
    To *Node         `json:"-"`
    // To *Node           `json:"-"`
    HoldingVal float64 `json:"holding"`
}

type Node struct {
    Value float64                      `json:"value"`  // should this be float32?  idk
    OutgoingConnections []*Connection  `json:"-"`  //which nodes to read from
    IncomingConnections []*Connection  `json:"-"`  //which nodes to read from
    Position [3]int                    `json:"position"`
}

type Network struct {
    // InputNodes []*Node   `json:"inputNodes"`
    Nodes []*Node        `json:"nodes"`
    // OutputNodes []*Node  `json:"outputNodes"`
    // MaxLevel int         `json:"maxLevel"`
    // MaxCycles int        `json:"maxCycles"`
    // CurCycle int         `json:"curCycle"`
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
    // TODO - REWORK
    var final float64
    for _, conn := range n.IncomingConnections {
        final = final + conn.HoldingVal*conn.Strength
    }
    final = final / float64(len(n.IncomingConnections))
    n.Value = final
}

func (net *Network) Cycle() {
    // fake concurrency
    // first, set all the connections based on their nodes

    for _, node := range net.Nodes {
        for _, conn := range node.OutgoingConnections {
            conn.HoldingVal = node.Value
            conn.Strength = conn.Strength * node.Value
            if conn.Strength < 0.1 || conn.Strength > 10 { // TWEAK THIS MAX STRENGTH!
                conn.Strength = 0.1
            }
        }
        node.Value = 0
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
                newConn := &Connection{
                    To: potConNode,
                    // Strength: rand.Float64() + 0.5, // do random strength - from 0.5 to 1.5?
                    Strength: rand.Float64(), // do random strength
                    // Strength: 1,
                    // Strength: rand.Float64()*2, // do random strength - from 0 to 2 or something?
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

func MakeNetwork(dimensions [3]int) *Network {
    nodes := []*Node{}
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

func main() {
    start := time.Now()

    // [width, depth, height]
    NETWORK_SIZE := [3]int{25, 25, 25}
    myNet := MakeNetwork(NETWORK_SIZE)
    myNet.Connect()
    myNet.Stimulate([]Stimulus{
        Stimulus{
            Position: [3]int{1,1,1},
            Strength: 5,
        },
    })
    frames, err := strconv.Atoi(os.Args[1])
    if err != nil {
        fmt.Println(err)
        return
    }
    myNet.GenerateAnim(frames)

    elapsed := time.Since(start)
    fmt.Printf("Took %s\n", elapsed)
}