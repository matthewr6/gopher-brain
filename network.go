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

// should one axon/connection connect to multiple neurons that are close by?
// http://changingminds.org/explanations/brain/parts_brain/neuron.htm
// how would this work with terminals then?
type Connection struct {
    To []*Node         `json:"-"`
    HoldingVal int     `json:"holding"`
    Terminals int      `json:"terminals"`
    Excitatory bool    `json:"excitatory"`
}

type Node struct {
    Value int                          `json:"value"`
    OutgoingConnection *Connection     `json:"-"`  //which node to send to
    IncomingConnections []*Connection  `json:"-"`  //which nodes to read from
    Position [3]int                    `json:"position"`
}

type Network struct {
    Nodes []*Node           `json:"nodes"`
    // Sensors []*Sensor       `json:"sensors"` // todo - add to state.go
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
}

func RandFloat(min, max float64) float64 {
    randFloat := rand.Float64()
    diff := max - min
    r := randFloat * diff
    return min + r
}

// todo - somewhere in here update the sensors... probably right after all the nodes update
func (net *Network) Cycle() {
    // fake concurrency
    // first, set all the connections based on their nodes

    for _, node := range net.Nodes {
        node.OutgoingConnection.HoldingVal = node.Value
        node.Value = 0
    }

    // then set all the nodes based on connections
    for _, node := range net.Nodes {
        node.Update()
    }

    // also update nodes that receive sensory information
    // for _, sensor := range net.Sensors {
    //     sensor.Update()
    // }

    // then clear the connections
    // do I still need this? doubtful
    // for _, node := range net.Nodes {
    //     node.OutgoingConnection.HoldingVal = 0
    // }
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

func NodeExistsIn(node *Node, nodes []*Node) bool {
    for _, potNode := range nodes {
        if (node == potNode) {
            return true
        }
    }
    return false
}

func (net *Network) Connect() {
    for _, node := range net.Nodes {
        // get the closest nodes and select one randomly to connect to

        // TODO - PRIORITY - rework - get a single possible node and then find areas around it,
        // similar to the way a sensor "sphere" is set up
        // TODO - add to state.go

        possibleConnections := []*Node{}
        for _, potConNode := range net.Nodes {
            if ThreeDimDist(node.Position, potConNode.Position) < 1.75 && node != potConNode {
                // todo - if the node already has more than X incoming connections, don't append?
                possibleConnections = append(possibleConnections, potConNode)
            }
        }
        // select the X connections here
        numAxonTerminals := rand.Intn(4) + 1 // TODO - HOW MANY POSSIBLE "TO" NEURONS
        nodesToConnect := []*Node{}
        for i := 0; i < numAxonTerminals; i++ {
            potNode := possibleConnections[rand.Intn(len(possibleConnections))]
            if !NodeExistsIn(potNode, nodesToConnect) {
                nodesToConnect = append(nodesToConnect, potNode)
            }
        }
        numTerminals := rand.Intn(3) + 1 // TODO - HOW MANY POSSIBLE TERMINALS
        var excitatory bool
        // should this have a higher probability of being excitatory?
        if rand.Intn(3) != 0 {
            excitatory = true
        }
        newConn := &Connection{
            To: nodesToConnect,
            Terminals: numTerminals,
            Excitatory: excitatory,
        }
        node.OutgoingConnection = newConn
        for _, nodeToConnect := range nodesToConnect {
            nodeToConnect.IncomingConnections = append(nodeToConnect.IncomingConnections, newConn)
        }
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