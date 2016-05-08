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
    To *Node           `json:"-"`
    HoldingVal int     `json:"holding"`
    Terminals int      `json:"terminals"` // like strenth - ADD THIS TO STATE.GO
    Excitatory bool    `json:"excitatory"` // TODO - ADD THIS TO STATE.GO
}

type Node struct {
    Value int                          `json:"value"`
    OutgoingConnection *Connection     `json:"axon"`  //which node to send to
    IncomingConnections []*Connection  `json:"-"`  //which nodes to read from
    Position [3]int                    `json:"position"`
}

type Network struct {
    Nodes []*Node `json:"nodes"`
}

type Stimulus struct {
    Position [3]int   `json:"position"`
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
    // calculate if it's going to fire or not - calculate sum and then set to 1, 0
    // base sum on excitatory/inhibiting
    sum := 0

    for _, conn := range n.IncomingConnections {
        if conn.Excitatory {
            sum = sum + (conn.HoldingVal*conn.Terminals)
        } else {
            sum = sum - (conn.HoldingVal*conn.Terminals)
        }
    }

    if sum >= 1 { // do 1 for threshold?
        //things
        //accept the value
        //or else just stay at 0
        n.Value = 1
    }

    // var final float64
    // for _, conn := range n.IncomingConnections {
    //     final = final + conn.HoldingVal*conn.Strength
    // }
    // n.Value = final
}

func RandFloat(min, max float64) float64 {
    randFloat := rand.Float64()
    diff := max - min
    r := randFloat * diff
    return min + r
}

func (net *Network) Cycle() {
    // fake concurrency
    // first, set all the connections based on their nodes

    for _, node := range net.Nodes {
        node.OutgoingConnection.HoldingVal = node.Value
        // bother with strengths?
        node.Value = 0
    }

    // then set all the nodes based on connections

    for _, node := range net.Nodes {
        node.Update()
    }

    // then clear the connections
    // do I still need this? doubtful
    for _, node := range net.Nodes {
        node.OutgoingConnection.HoldingVal = 0
    }
}

func (n Node) String() string {
    jsonRep, _ := json.MarshalIndent(n, "", "    ")
    return string(jsonRep)
}


//rework - stimulate certain neurons, but don't bother with strength - just 1
func (net *Network) Stimulate(stimuli []Stimulus) {
    for _, stim := range stimuli {
        var applyTo *Node;
        for _, node := range net.Nodes {
            if node.Position == stim.Position {
                applyTo = node
                break
            }
        }
        applyTo.Value = 1
    }
}

// is this still needed?
// for cool viz, sure
func (net *Network) RandomizeValues() {
    for _, node := range net.Nodes {
        temp := rand.Float32()
        if temp < 0.5 {
            node.Value = 0
        } else {
            node.Value = 1
        }
    }
}

func ThreeDimDist(p1, p2 [3]int) float64 {
    ans := (p1[0]-p2[0])*(p1[0]-p2[0]) + (p1[1]-p2[1])*(p1[1]-p2[1]) + (p1[2]-p2[2])*(p1[2]-p2[2])
    return math.Sqrt(float64(ans))
}

func (net *Network) Connect() {
    for _, node := range net.Nodes {
        // get the closest nodes and select one randomly to connect to
        possibleConnections := []*Node{}
        for _, potConNode := range net.Nodes {
            if ThreeDimDist(node.Position, potConNode.Position) < 1.75 && node != potConNode {
                possibleConnections = append(possibleConnections, potConNode)
            }
        }
        // select the one connection here
        nodeToConnect := possibleConnections[rand.Intn(len(possibleConnections))]
        numTerminals := rand.Intn(3) + 1 // TODO - HOW MANY POSSIBLE TERMINALS
        var excitatory bool
        randTest := rand.Float32()
        if randTest < 0.5 {
            excitatory = true
        }
        newConn := &Connection{
            To: nodeToConnect,
            Terminals: numTerminals,
            Excitatory: excitatory,
        }
        node.OutgoingConnection = newConn
        nodeToConnect.IncomingConnections = append(nodeToConnect.IncomingConnections, newConn)
    }
}

func (net Network) String() string {
    jsonRep, _ := json.MarshalIndent(net, "", "    ")
    return string(jsonRep)
}

func (net Network) DumpJSON(name string) {
    f, _ := os.Create(fmt.Sprintf("./frames/net_%v.json", name))
    f.WriteString(net.String())
    f.Close()
}

func MakeNetwork(dimensions [3]int, blank bool) *Network {
    nodes := []*Node{}
    math.Mod(5, 5)
    for i := 1; i <= dimensions[0]; i++ {
        for j := 1; j <= dimensions[1]; j++ {
            for k := 1; k <= dimensions[2]; k++ {
                var newValue int
                var randTest float32
                if !blank {
                    randTest = rand.Float32()
                    if randTest < 0.5 {
                        newValue = 0
                    } else {
                        newValue = 1
                    }
                }
                nodes = append(nodes, &Node{
                    Value: newValue,
                    Position: [3]int{i, j, k},
                    IncomingConnections: []*Connection{},
                })
            }
        }
    }
    return &Network {
        Nodes: nodes,
    }
}

func (net *Network) GenerateAnim(frames int) {
    os.Mkdir("frames", 755)
    for frame := 0; frame < frames; frame++ {
        net.DumpJSON(strconv.Itoa(frame))
        net.Cycle()
    }
}