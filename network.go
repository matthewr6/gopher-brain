package main

import (
    "fmt"
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
    Value int                          `json:"value"` // todo- maybe use a bool of true or false, and rename this "firing"?
    OutgoingConnection *Connection     `json:"-"`  //which node to send to
    IncomingConnections []*Connection  `json:"-"`  //which nodes to read from
    Position [3]int                    `json:"position"`
    Id string                          `json:"id"` //"x|y|z"
}

// let's say y=0 is the front of the "brain"
type Network struct {
    Nodes [][][]*Node           `json:"nodes"`
    LeftHemisphere [][][]*Node  `json:"leftHemisphere"`
    RightHemisphere [][][]*Node  `json:"rightHemisphere"`
    Dimensions [3]int           `json:"-"`
    Sensors []*Sensor           `json:"-"`
    Outputs []*Output           `json:"-"`
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
        // let's just wrap it in this for now...
        if conn.To[n] != nil {
            // fmt.Println()
            // fmt.Println(n.Id)
            // fmt.Println(conn.To[n])
            if conn.To[n].Excitatory {
                sum = sum + (float64(conn.HoldingVal) * conn.To[n].Strength)
            } else {
                sum = sum - (float64(conn.HoldingVal) * conn.To[n].Strength)
            }

            // reassess connections here
            // todo - calculate how much to increase/decrease connection strength by
            // https://www.reddit.com/r/askscience/comments/1bb5br/what_physically_happens_when_neural_connections/
            if conn.HoldingVal == 0 {
                // the previous node *didn't* fire
                conn.To[n].Strength -= 0.05
            } else {
                // the previous node *did* fire
                conn.To[n].Strength += 0.05
            }

            // todo - thresholds
            if conn.To[n].Strength > 2.25 {
                // max strength?
                conn.To[n].Strength = 2.25
            }
            // the below has to be at the end
            // it's not a pretty way to resolve it but it works
            // maybe use `continue`

            // todo - this isn't working
            // oh have to delete conn from incomingconnections
            // maybe this stuff in new loop?
            if conn.To[n].Strength < 0.25 {
                // remove?  different threshold?
                delete(conn.To, n)
                // fmt.Println("deleting")
                // fmt.Println(conn.To[n])
            }
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
    numPossible := rand.Intn(15 - 5) + 5
    stDev := 3.0 // todo, what number?
    for i := 0; i < numPossible; i++ {
        // todo - wrapper function for this since it's used so much
        potX := int(rand.NormFloat64() * stDev) + center[0]
        potY := int(rand.NormFloat64() * stDev) + center[1]
        potZ := int(rand.NormFloat64() * stDev) + center[2]
        for potX < 0 || potX >= net.Dimensions[0] {
            potX = int(rand.NormFloat64() * stDev) + center[0]
        }
        for potY < 0 || potY >= net.Dimensions[0] {
            potY = int(rand.NormFloat64() * stDev) + center[1]
        }
        for potZ < 0 || potZ >= net.Dimensions[0] {
            potZ = int(rand.NormFloat64() * stDev) + center[2]
        }
        potCenter := [3]int{potX, potY, potZ}
        possibleExtensions = append(possibleExtensions, net.FindNode(potCenter))
    }
    // could merge this into the above loop...
    for _, potNode := range possibleExtensions {
        _, exists := node.OutgoingConnection.To[potNode]
        if potNode.Value != 0 && !exists { // todo doesthis check work?
            excitatory := false
            if rand.Intn(2) != 0 {
                excitatory = true
            }
            node.OutgoingConnection.To[potNode] = &ConnInfo{
                Strength: RandFloat(0.50, 1.50),
                Excitatory: excitatory,
            }
        } 
    }
}

// let's see which one causes the most overhead...
// or it might just be all of them
func (net *Network) Cycle() {
    // fake concurrency
    // first, set all the connections based on their nodes

    net.ForEachNode(func(node *Node, pos [3]int) {
        // todo - search for nodes to connect to?
        // what should this be on, the node or the connection?
        // also make sure the order is good
        if node.Value != 0 {
            net.AddConnections(node)
        }
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

func NodeExistsIn(node *Node, nodes []*Node) bool {
    for _, potNode := range nodes {
        if (node == potNode) {
            return true
        }
    }
    return false
}

// todo
// call this after connecting first hemisphere
func (net *Network) Mirror() {
    // invert in x direction
    leftHemisphere := [][][]*Node{}
    for i := len(net.RightHemisphere)-1; i >= 0; i-- {
        leftHemisphere = append(leftHemisphere, net.RightHemisphere[i])
    }
    net.LeftHemisphere = leftHemisphere
    net.ForEachRightHemisphereNode(func(node *Node, pos [3]int) {
        actualNode := net.FindLeftHemisphereNode(pos)

        var newConnection = &Connection{
            HoldingVal: 0,
        }

        newConnection.To = make(map[*Node]*ConnInfo)
        newConnection.Center[0] = (net.Dimensions[0]-1) - newConnection.Center[0]
        // now have to do the "to" stuff
        centralConnNode := net.FindLeftHemisphereNode(newConnection.Center) // is this correct

        // start redundancy
        numAxonTerminals := rand.Intn(3) + 1 // todo
        nodesToConnect := []*Node{
            centralConnNode,
        }
        stDev := 1.5
        for i := 0; i < numAxonTerminals; i++ {
            potPos := [3]int{-1, -1, -1}
            for potPos[0] < 0 || potPos[1] < 0 || potPos[2] < 0 || potPos[0] >= net.Dimensions[0] || potPos[1] >= net.Dimensions[1] || potPos[2] >= net.Dimensions[2] {
                potPos = [3]int{
                    int(rand.NormFloat64() * stDev) + centralConnNode.Position[0],
                    int(rand.NormFloat64() * stDev) + centralConnNode.Position[1],
                    int(rand.NormFloat64() * stDev) + centralConnNode.Position[2],
                }
            }
            potNode := net.FindLeftHemisphereNode(potPos)
            // potNode := possibleConnections[rand.Intn(len(possibleConnections))]
            if !NodeExistsIn(potNode, nodesToConnect) && potNode != actualNode {
                nodesToConnect = append(nodesToConnect, potNode)
            }
        }
        var excitatory bool
        toNodes := make(map[*Node]*ConnInfo)
        for _, connNode := range nodesToConnect {
            // should this have a higher probability of being excitatory?
            if rand.Intn(2) != 0 {
                excitatory = true
            }
            toNodes[connNode] = &ConnInfo{
                Strength: RandFloat(0.75, 1.75),
                Excitatory: excitatory,
            }
        }
        newConn := &Connection{
            To: toNodes,
            Center: centralConnNode.Position,
        }
        actualNode.OutgoingConnection = newConn
        for _, nodeToConnect := range nodesToConnect {
            nodeToConnect.IncomingConnections = append(nodeToConnect.IncomingConnections, newConn)
        }
        // end redundancy
        // maybe abstract some stuff in the Connect() function into another function?
    })
    // todo - somehow concatenate both into the Nodes attribute
    nodes := net.RightHemisphere // pointer issues?
    for _, nodePlane := range net.LeftHemisphere {
        nodes = append(nodes, nodePlane)
    }
    net.Nodes = nodes
    // reset ids
    net.ForEachNode(func(node *Node, pos [3]int) {
        node.Id = fmt.Sprintf("%v|%v|%v", pos[0], pos[1], pos[2])
        node.Position = pos
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
            for potY < 0 || potY >= net.Dimensions[0] {
                potY = int(rand.NormFloat64() * stDev) + node.Position[1]
            }
            for potZ < 0 || potZ >= net.Dimensions[0] {
                potZ = int(rand.NormFloat64() * stDev) + node.Position[2]
            }
            center = [3]int{potX, potY, potZ}
        }
        centralConnNode := net.FindRightHemisphereNode(center)

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
            potNode := net.FindRightHemisphereNode(potPos)
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
            Center: centralConnNode.Position,
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
        RightHemisphere: nodes,
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