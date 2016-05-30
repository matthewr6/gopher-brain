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

// http://www.scientificamerican.com/article/ask-the-brains-aug-08/
// maybe do a map of a connected node to its strength and/or its excitatory/inhibitory

// todo - state.go
// array with {node: node, strength: strength}
// unique ids for each node, maybe?  based on position?  like "x|y|z" concatenated?

type ConnInfo struct {
    Excitatory bool  `json:"excitatory"`
    Strength float64 `json:"strength"`
}

type Connection struct {
    To map[*Node]*ConnInfo  `json:"to"`
    HoldingVal int          `json:"holding"`
}

type Node struct {
    Value int                          `json:"value"`
    OutgoingConnection *Connection     `json:"-"`  //which node to send to
    IncomingConnections []*Connection  `json:"-"`  //which nodes to read from
    Position [3]int                    `json:"position"`
    Id string                          `json:"id"` //"x|y|z"
}

type Network struct {
    Nodes [][][]*Node       `json:"nodes"`
    Dimensions [3]int       `json:"-"`
    Sensors []*Sensor       `json:"-"`
    Outputs []*Output       `json:"-"`
}

func (c Connection) String() string {
    jsonRep, _ := json.MarshalIndent(c, "", "    ")
    return string(jsonRep)
}

func (n *Node) Update() {
    // calculate if it's going to fire or not - calculate sum and then set to 1, 0
    // base sum on excitatory/inhibiting
    var sum float64

    for _, conn := range n.IncomingConnections {
        if conn.To[n].Excitatory {
            sum = sum + (float64(conn.HoldingVal) * conn.To[n].Strength)
        } else {
            sum = sum - (float64(conn.HoldingVal) * conn.To[n].Strength)
        }
    }

    if sum >= 1.0 { // do 1 for threshold?
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

// let's see which one causes the most overhead...
// or it might just be all of them
func (net *Network) Cycle() {
    // fake concurrency
    // first, set all the connections based on their nodes

    net.ForEachNode(func(node *Node, pos [3]int) {
        node.OutgoingConnection.HoldingVal = node.Value
        node.Value = 0
    })

    // then set all the nodes based on connections
    net.ForEachNode(func(node *Node, pos [3]int) {
        node.Update()
    })

    // also update nodes that receive sensory information
    for _, sensor := range net.Sensors {
        sensor.In(sensor.Nodes, sensor.Stimulated)
    }

    for _, output := range net.Outputs {
        output.Value = output.Out(output.Nodes)
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

// is this still needed?
// for cool viz, sure
func (net *Network) RandomizeValues() {
    net.ForEachNode(func(node *Node, pos [3]int) {
        temp := rand.Float32()
        if temp < 0.5 {
            node.Value = 0
        } else {
            node.Value = 1
        }
    })
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
    net.ForEachNode(func(node *Node, pos [3]int) {
        // get the closest nodes and select one randomly to connect to
        stDev := 3.0 // what should it be?
        center := node.Position
        for center == node.Position {
            potX := int(rand.NormFloat64() * stDev) + node.Position[0]
            potY := int(rand.NormFloat64() * stDev) + node.Position[1]
            potZ := int(rand.NormFloat64() * stDev) + node.Position[2]
            for potX < 0 || potX >= net.Dimensions[0] {
                potX = int(rand.NormFloat64() * stDev) + node.Position[0]
            }
            for potY < 0 || potY >= net.Dimensions[0] {
                potY = int(rand.NormFloat64() * stDev) + node.Position[1]
            }
            for potZ < 0 || potZ >= net.Dimensions[0] {
                potZ = int(rand.NormFloat64() * stDev) + node.Position[2]
            }
            center = [3]int{potX, potY, potZ}
        }
        centralConnNode := FindNode(center, net.Nodes)

        // select the X connections here
        numAxonTerminals := rand.Intn(3) + 1 // TODO - HOW MANY POSSIBLE "TO" NEURONS - 3 max seems good
        nodesToConnect := []*Node{
            centralConnNode,
        }
        stDev = 1.5
        for i := 0; i < numAxonTerminals; i++ {
            potPos := [3]int{-1, -1, -1}
            for potPos[0] < 0 || potPos[1] < 0 || potPos[2] < 0 || potPos[0] >= net.Dimensions[0] || potPos[1] >= net.Dimensions[1] || potPos[2] >= net.Dimensions[2] {
                potPos = [3]int{
                    int(rand.NormFloat64() * stDev) + centralConnNode.Position[0],
                    int(rand.NormFloat64() * stDev) + centralConnNode.Position[1],
                    int(rand.NormFloat64() * stDev) + centralConnNode.Position[2],
                }
            }
            potNode := FindNode(potPos, net.Nodes)
            // potNode := possibleConnections[rand.Intn(len(possibleConnections))]
            if !NodeExistsIn(potNode, nodesToConnect) && potNode != node {
                nodesToConnect = append(nodesToConnect, potNode)
            }
        }

        var excitatory bool
        toNodes := make(map[*Node]*ConnInfo)
        for _, node := range nodesToConnect {
            // should this have a higher probability of being excitatory?
            if rand.Intn(2) != 0 {
                excitatory = true
            }
            toNodes[node] = &ConnInfo{
                Strength: RandFloat(0.75, 1.75),
                Excitatory: excitatory,
            }
        }

        newConn := &Connection{
            To: toNodes,
        }
        node.OutgoingConnection = newConn
        for _, nodeToConnect := range nodesToConnect {
            nodeToConnect.IncomingConnections = append(nodeToConnect.IncomingConnections, newConn)
        }
    })
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

// probably isn't best way to do this - will have to rethink
func MakeNetwork(dimensions [3]int, blank bool) *Network {
    nodes := [][][]*Node{}
    for i := 0; i < dimensions[0]; i++ {
        iDim := [][]*Node{}
        for j := 0; j < dimensions[1]; j++ {
            jDim := []*Node{}
            for k := 0; k < dimensions[2]; k++ {
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
                jDim = append(jDim, &Node{
                    Value: newValue,
                    Position: [3]int{i, j, k},
                    IncomingConnections: []*Connection{},
                    Id: fmt.Sprintf("%v|%v|%v", i, j, k),
                })
            }
            iDim = append(iDim, jDim)
        }
        nodes = append(nodes, iDim)
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
    fmt.Print("\n-----SENSORS-----\n")
    for _, sensor := range net.Sensors {
        active := "x"
        if sensor.Stimulated {
            active = "o"
        }
        fmt.Printf("%v: %v\n", sensor.Name, active)
    }
    fmt.Print("\n-----OUTPUTS-----\n")
    for _, output := range net.Outputs {
        fmt.Printf("%v: %v\n", output.Name, output.Value)
    }
    term.HideCursor()
}

func (net *Network) ForEachNode(handler func(*Node, [3]int)) {
    for i := range net.Nodes {
        for j := range net.Nodes[i] {
            for k := range net.Nodes[i][j] {
                handler(net.Nodes[i][j][k], [3]int{i, j, k})
            }
        }
    }
}