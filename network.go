package brain

import (
    "fmt"
    "math/rand"
    "os"
    "strconv"
    "encoding/json"
)

/*
    & is "address of"
    * is "value at address"
*/

// http://www.scientificamerican.com/article/ask-the-brains-aug-08/

type ConnInfo struct {
    Excitatory bool  `json:"excitatory"`
    Strength float64 `json:"strength"`
}

type Connection struct {
    To map[*Node]*ConnInfo  `json:"to"`
    HoldingVal int          `json:"holding"`
    Center [3]int           `json:"center"`
}

type Node struct {
    Value int                                  `json:"value"`
    OutgoingConnection *Connection             `json:"-"`  //which node to send to
    IncomingConnections map[*Node]*Connection  `json:"-"`  //which nodes to read from
    Position [3]int                            `json:"position"`
    Id string                                  `json:"id"` //"x|y|z"
}

// let's say y=0 is the front of the "brain"
type Network struct {
    Nodes [][][]*Node             `json:"nodes"`
    LeftHemisphere [][][]*Node    `json:"-"`
    RightHemisphere [][][]*Node   `json:"-"`
    Dimensions [3]int             `json:"-"`
    Sensors map[string]*Sensor    `json:"sensors"`
    Outputs map[string]*Output    `json:"outputs"`
}

func (c Connection) String() string {
    jsonRep, _ := json.MarshalIndent(c, "", "    ")
    return string(jsonRep)
}

func (n *Node) Update() {
    // calculate if it's going to fire or not - calculate sum and then set to 1, 0
    // base sum on excitatory/inhibiting
    var sum float64

    for from, conn := range n.IncomingConnections {
        // let's just wrap it in this for now...
        if conn.To[n].Excitatory {
            sum = sum + (float64(conn.HoldingVal) * conn.To[n].Strength)
        } else {
            sum = sum - (float64(conn.HoldingVal) * conn.To[n].Strength)
        }

        // reassess connections here
        // todo - MAGIC # - calculate how much to increase/decrease connection strength by
        // https://www.reddit.com/r/askscience/comments/1bb5br/what_physically_happens_when_neural_connections/
        if conn.HoldingVal == 0 {
            // the previous node *didn't* fire
            conn.To[n].Strength -= 0.05
        } else {
            // the previous node *did* fire
            conn.To[n].Strength += 0.05
        }

        // todo - MAGIC # - thresholds
        if conn.To[n].Strength > 2.25 {
            // max strength?
            conn.To[n].Strength = 2.25
        }
        if conn.To[n].Strength < 0.25 {
            // different threshold?
            delete(conn.To, n)
            delete(n.IncomingConnections, from)
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

func (net *Network) AddConnections(node *Node) {
    center := node.OutgoingConnection.Center
    possibleExtensions := []*Node{}
    numPossible := rand.Intn(15 - 5) + 5 // 10 to 15
    // todo - MAGIC # - what number?
    stDev := 3.0
    for i := 0; i < numPossible; i++ {
        // todo - POSSIBLE - wrapper function for this since it's used so much
        potX := int(rand.NormFloat64() * stDev) + center[0]
        potY := int(rand.NormFloat64() * stDev) + center[1]
        potZ := int(rand.NormFloat64() * stDev) + center[2]
        for potX < 0 || potX >= (net.Dimensions[0] * 2) {
            potX = int(rand.NormFloat64() * stDev) + center[0]
        }
        for potY < 0 || potY >= net.Dimensions[1] {
            potY = int(rand.NormFloat64() * stDev) + center[1]
        }
        for potZ < 0 || potZ >= net.Dimensions[2] {
            potZ = int(rand.NormFloat64() * stDev) + center[2]
        }
        potCenter := [3]int{potX, potY, potZ}
        possibleExtensions = append(possibleExtensions, net.FindNode(potCenter))
    }
    // todo - POSSIBLE
    // could merge this into the above loop...
    for _, potNode := range possibleExtensions {
        _, exists := node.OutgoingConnection.To[potNode]
        if potNode.Value != 0 && !exists {
            excitatory := false
            if rand.Intn(2) != 0 {
                excitatory = true
            }
            node.OutgoingConnection.To[potNode] = &ConnInfo{
                Strength: RandFloat(0.50, 1.50),
                Excitatory: excitatory,
            }
            potNode.IncomingConnections[node] = node.OutgoingConnection
        } 
    }
}

func (node *Node) UpdateOutgoingCenter() {
    x := 0
    y := 0
    z := 0
    numOutgoing := len(node.OutgoingConnection.To)
    if numOutgoing < 0 {
        for to := range node.OutgoingConnection.To {
            x += to.Position[0]
            y += to.Position[1]
            z += to.Position[2]
        }
        x = x / numOutgoing
        y = y / numOutgoing
        z = z / numOutgoing
        node.OutgoingConnection.Center = [3]int{x, y, z}
    }
}

// let's see which one causes the most overhead...
// or it might just be all of them
func (net *Network) Cycle() {
    // fake concurrency
    // first, set all the connections based on their nodes

    net.ForEachNode(func(node *Node, pos [3]int) {
        if node.Value != 0 {
            net.AddConnections(node)
        }
        node.OutgoingConnection.HoldingVal = node.Value
        node.Value = 0
    })
    

    // then set all the nodes based on connections
    net.ForEachNode(func(node *Node, pos [3]int) {
        node.Update()
        node.UpdateOutgoingCenter()
    })


    // also update nodes that receive sensory information
    for _, output := range net.Outputs {
        output.Value = output.Out(output.Nodes)
    }

    for _, sensor := range net.Sensors {
        sensor.In(sensor.Nodes, sensor.Influences)
    }
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

func NodeExistsIn(node *Node, nodes []*Node) bool {
    for _, potNode := range nodes {
        if (node == potNode) {
            return true
        }
    }
    return false
}

func (net *Network) ConnectHemispheres() {
    net.ForEachNode(func(node *Node, pos [3]int) {
        centralConnNode := net.FindNode(node.OutgoingConnection.Center)
        // select the X connections here
        // TODO - MAGIC # - HOW MANY POSSIBLE "TO" NEURONS - 3 max seems good
        numAxonTerminals := rand.Intn(3) + 1
        nodesToConnect := []*Node{
            centralConnNode,
        }
        stDev := 1.5
        for i := 0; i < numAxonTerminals; i++ {
            potPos := [3]int{-1, -1, -1}
            for potPos[0] < 0 || potPos[1] < 0 || potPos[2] < 0 || potPos[0] >= net.Dimensions[0]*2 || potPos[1] >= net.Dimensions[1] || potPos[2] >= net.Dimensions[2] {
                potPos = [3]int{
                    int(rand.NormFloat64() * stDev) + centralConnNode.Position[0],
                    int(rand.NormFloat64() * stDev) + centralConnNode.Position[1],
                    int(rand.NormFloat64() * stDev) + centralConnNode.Position[2],
                }
            }
            potNode := net.FindNode(potPos)
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

        node.OutgoingConnection.To = toNodes
        for _, nodeToConnect := range nodesToConnect {
            nodeToConnect.IncomingConnections[node] = node.OutgoingConnection
        }
    })
}

func (net *Network) Mirror() {
    // invert in x direction
    leftHemisphere := [][][]*Node{}
    for i := len(net.RightHemisphere)-1; i >= 0; i-- {
        // POINTER CRAPS - NODES
        rightPlane := [][]*Node{}
        for _, rightRow := range net.RightHemisphere[i] {
            aRightRow := []*Node{}
            for _, rightNode := range rightRow {
                newNode := &Node{}
                *newNode = *rightNode
                newNode.IncomingConnections = make(map[*Node]*Connection)
                aRightRow = append(aRightRow, newNode)
            }
            rightPlane = append(rightPlane, aRightRow)
        }
        // leftHemisphere = append(leftHemisphere, net.RightHemisphere[i])
        leftHemisphere = append(leftHemisphere, rightPlane)
    }
    net.LeftHemisphere = leftHemisphere
    net.ForEachRightHemisphereNode(func(node *Node, pos [3]int) {
        actualNode := net.FindLeftHemisphereNode(pos)

        newCenter := node.OutgoingConnection.Center
        newCenter[0] = (net.Dimensions[0]-1) - node.OutgoingConnection.Center[0]

        newCenter[0] += net.Dimensions[0]
        newConn := &Connection{
            Center: newCenter,
        }
        actualNode.OutgoingConnection = newConn
    })
    net.Nodes = append(net.RightHemisphere, net.LeftHemisphere...)
    net.ForEachNode(func(node *Node, pos [3]int) {
        node.Position = pos
        node.Id = fmt.Sprintf("%v|%v|%v", pos[0], pos[1], pos[2])
    })
}

func (net *Network) Connect() {
    net.ForEachRightHemisphereNode(func(node *Node, pos [3]int) {
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
            for potY < 0 || potY >= net.Dimensions[1] {
                potY = int(rand.NormFloat64() * stDev) + node.Position[1]
            }
            for potZ < 0 || potZ >= net.Dimensions[2] {
                potZ = int(rand.NormFloat64() * stDev) + node.Position[2]
            }
            center = [3]int{potX, potY, potZ}
        }

        newConn := &Connection{
            Center: center,
        }
        node.OutgoingConnection = newConn
    })
}

func (net Network) String() string {
    jsonRep, _ := json.MarshalIndent(net, "", "    ")
    return string(jsonRep)
}

func (net Network) DumpJSON(name string) {
    f, _ := os.Create(fmt.Sprintf("%v/frames/net_%v.json", directory, name))
    f.WriteString(net.String())
    f.Close()
}

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
                        newValue = 0//1
                    }
                }
                jDim = append(jDim, &Node{
                    Value: newValue,
                    Position: [3]int{i, j, k},
                    IncomingConnections: make(map[*Node]*Connection),
                    Id: fmt.Sprintf("%v|%v|%v", i, j, k),
                })
            }
            iDim = append(iDim, jDim)
        }
        nodes = append(nodes, iDim)
    }

    return &Network {
        Dimensions: dimensions,
        RightHemisphere: nodes,
        Sensors: make(map[string]*Sensor),
        Outputs: make(map[string]*Output),
    }
}

func (net *Network) GenerateAnim(frames int) {
    os.Mkdir("frames", 755)
    for frame := 0; frame < frames; frame++ {
        net.DumpJSON(strconv.Itoa(frame))
        net.Cycle()
    }
}

func (net *Network) AnimateUntilDone() {
    os.Mkdir("frames", 755)
    frame := 0
    for running {
        frameStr := strconv.Itoa(frame)
        net.DumpJSON(frameStr)
        net.Cycle()
        frame++
    }
}