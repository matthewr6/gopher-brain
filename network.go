package main

import (
    "fmt"
    "math"
    "math/rand"
    "os"
    "strconv"
    "encoding/json"
    "time"

    "github.com/jteeuwen/keyboard"
    term "github.com/nsf/termbox-go"
)

/*
    & is "address of"
    * is "value at address"
*/

// should one axon/connection connect to multiple neurons that are close by?
// http://changingminds.org/explanations/brain/parts_brain/neuron.htm
//http://cogsci.stackexchange.com/questions/9144/how-many-dendrite-connections-vs-axon-terminals-does-a-multipolar-cerebral-neuro
// todo how would this work with terminals then?
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
    Dimensions [3]int       `json:"dimensions"`
    Sensors []*Sensor       `json:"sensors"`
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
    for _, sensor := range net.Sensors {
        sensor.Update()
    }

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
        stDev := 3.0 // what should it be?
        center := node.Position
        for center == node.Position {
            potX := int(rand.NormFloat64() * stDev) + node.Position[0]
            potY := int(rand.NormFloat64() * stDev) + node.Position[1]
            potZ := int(rand.NormFloat64() * stDev) + node.Position[2]
            for potX <= 0 || potX > net.Dimensions[0] {
                potX = int(rand.NormFloat64() * stDev) + node.Position[0]
            }
            for potY <= 0 || potY > net.Dimensions[0] {
                potY = int(rand.NormFloat64() * stDev) + node.Position[1]
            }
            for potZ <= 0 || potZ > net.Dimensions[0] {
                potZ = int(rand.NormFloat64() * stDev) + node.Position[2]
            }
            center = [3]int{potX, potY, potZ}
        }
        centralConnNode := FindNode(center, net.Nodes)

        // select the X connections here
        numAxonTerminals := rand.Intn(3) + 1 // TODO - HOW MANY POSSIBLE "TO" NEURONS?
        nodesToConnect := []*Node{
            centralConnNode,
        }
        stDev = 1.5
        for i := 0; i < numAxonTerminals; i++ {
            potPos := [3]int{
                int(rand.NormFloat64() * stDev) + centralConnNode.Position[0],
                int(rand.NormFloat64() * stDev) + centralConnNode.Position[1],
                int(rand.NormFloat64() * stDev) + centralConnNode.Position[2],
            }
            potNode := FindNode(potPos, net.Nodes)
            // potNode := possibleConnections[rand.Intn(len(possibleConnections))]
            if !NodeExistsIn(potNode, nodesToConnect) && potNode != node {
                nodesToConnect = append(nodesToConnect, potNode)
            }
        }

        // do I even want this now?
        numTerminals := rand.Intn(2) + 1 // TODO - HOW MANY POSSIBLE TERMINALS

        var excitatory bool
        // should this have a higher probability of being excitatory?
        if rand.Intn(2) != 0 {
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
        Dimensions: dimensions,
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

func (net *Network) AnimateUntilDone(ms int) {
    os.Mkdir("frames", 755)
    frame := 0
    for running {
        time.Sleep(time.Duration(ms) * time.Millisecond)
        frameStr := strconv.Itoa(frame)
        net.DumpJSON(frameStr)
        net.Cycle()
        // should print everything on one line, just because it's simpler
        // fmt.Print("\r" + frameStr)
        net.Info(frame)
        frame++
    }
}

func KeyboardPoll(kb keyboard.Keyboard) {
    for running {
        kb.Poll(term.PollEvent())
    }
}

func (net Network) Info(frame int) {
    term.SetCursor(0, 0)
    fmt.Printf("Frame %v\n", frame)
    for _, sensor := range net.Sensors {
        active := "inactive"
        if sensor.Stimulated {
            active = "active"
        }
        fmt.Printf("%v: %v\n", sensor.Name, active)
    }
    term.HideCursor()
}