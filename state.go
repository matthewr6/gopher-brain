package main

import (
    "fmt"
    "os"
    "encoding/json"
)

/*
    & is "address of"
    * is "value at address"
*/

// still have to reconcile the updated network.go stuff

type DisplayNetwork struct {
    Nodes []*DisplayNode     `json:"nodes"`
    Dimensions [3]int        `json:"dimensions"`
    Sensors []*DisplaySensor `json:"sensors"`
    // Connections []*DisplayConnection `json:"connections"`
}

type DisplayNode struct {
    Value int                               `json:"value"`
    Position [3]int                         `json:"position"`
    OutgoingConnection *DisplayConnection   `json:"axon"`
}

type DisplayConnection struct {
    To [][3]int           `json:"to"`
    HoldingVal int        `json:"holdingVal"`
    Terminals int         `json:"terminals"`
    Excitatory bool       `json:"excitatory"`
}

type DisplaySensor struct {
    Nodes[][3]int     `json:"nodes"`
    Excitatory bool   `json:"excitatory"`
    Trigger string    `json:"trigger"`
}

func (d DisplayNetwork) String() string {
    jsonRep, _ := json.MarshalIndent(d, "", "    ")
    return string(jsonRep)
}

// oh sweet jesus MORE INEFFICIENCY
func FindNode(position [3]int, potentialNodes []*Node) *Node {
    for _, potNode := range potentialNodes {
        if potNode.Position == position {
            return potNode
        }
    }
    return &Node{}
}

func LoadState(name string) *Network {
    fmt.Println("loading")
    datafile, err := os.Open(fmt.Sprintf("./state/%v_state.json", name))
    if err != nil {
        fmt.Println(err)
    }
    decoder := json.NewDecoder(datafile)
    importedNet := &DisplayNetwork{}
    decoder.Decode(&importedNet)
    datafile.Close()

    net := &Network{
        Nodes: []*Node{},
        Dimensions: importedNet.Dimensions,
    }
    // set nodes
    // this looks good
    for _, importedNode := range importedNet.Nodes {
        newNode := &Node{
            Value: importedNode.Value,
            Position: importedNode.Position,
            IncomingConnections: []*Connection{},
        }
        net.Nodes = append(net.Nodes, newNode)
    }
    // set connections
    // this part is super inefficient
    // still should optimize
    for _,  importedNode := range importedNet.Nodes {
        newConn := &Connection{
            HoldingVal: importedNode.OutgoingConnection.HoldingVal,
            Terminals: importedNode.OutgoingConnection.Terminals,
            Excitatory: importedNode.OutgoingConnection.Excitatory,
        }
        node := FindNode(importedNode.Position, net.Nodes)
        nodesToConnect := []*Node{}
        for _, nodePos := range importedNode.OutgoingConnection.To {
            // these similar names are gonna kill me
            nodeToConnect := FindNode(nodePos, net.Nodes)
            nodesToConnect = append(nodesToConnect, nodeToConnect)
            nodeToConnect.IncomingConnections = append(nodeToConnect.IncomingConnections, newConn)
        }
        node.OutgoingConnection = newConn
        newConn.To = nodesToConnect
    }

    // set sensors
    // this is also inefficient
    for _, importedSensor := range importedNet.Sensors {
        nodes := []*Node{}
        for _, nodePos := range importedSensor.Nodes {
            nodes = append(nodes, FindNode(nodePos, net.Nodes))
        }
        net.Sensors = append(net.Sensors, &Sensor{
            Nodes: nodes,
            Excitatory: importedSensor.Excitatory,
            Trigger: importedSensor.Trigger,
            Stimulated: false,
        })
    }
    return net
}

func (net Network) SaveState(name string) {
    fmt.Println("saving")
    os.Mkdir("state", 755)
    dispNet := DisplayNetwork{
        Nodes: []*DisplayNode{},
        Sensors: []*DisplaySensor{},
        Dimensions: net.Dimensions,
    }
    for _, sensor := range net.Sensors {
        positions := [][3]int{}
        for _, sensoryNode := range sensor.Nodes {
            positions = append(positions, sensoryNode.Position)
        }
        dispNet.Sensors = append(dispNet.Sensors, &DisplaySensor{
            Nodes: positions,
            Excitatory: sensor.Excitatory,
            Trigger: sensor.Trigger,
        })
    }
    for _, node := range net.Nodes {
        toPositions := [][3]int{}
        for _, outNode := range node.OutgoingConnection.To {
            toPositions = append(toPositions, outNode.Position)
        }
        dispConn := &DisplayConnection{
            To: toPositions,
            HoldingVal: node.OutgoingConnection.HoldingVal,
            Terminals: node.OutgoingConnection.Terminals,
            Excitatory: node.OutgoingConnection.Excitatory,
        }

        dispNode := &DisplayNode{
            Value: node.Value,
            Position: node.Position,
            OutgoingConnection: dispConn,
        }
        dispNet.Nodes = append(dispNet.Nodes, dispNode)
    }
    f, _ := os.Create(fmt.Sprintf("./state/%v_state.json", name))
    f.WriteString(dispNet.String())
    f.Close()
}