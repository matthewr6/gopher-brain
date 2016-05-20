package main

import (
    "fmt"
    "os"
    "encoding/json"

    "github.com/jteeuwen/keyboard"
)

/*
    & is "address of"
    * is "value at address"
*/

// still have to reconcile the updated network.go stuff

type DisplayNetwork struct {
    Nodes [][][]*DisplayNode   `json:"nodes"`
    Dimensions [3]int          `json:"dimensions"`
    Sensors []*DisplaySensor   `json:"sensors"`
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
    Name string       `json:"name"`
}

func (d DisplayNetwork) String() string {
    jsonRep, _ := json.MarshalIndent(d, "", "    ")
    return string(jsonRep)
}

// oh sweet jesus MORE INEFFICIENCY
// todo replace all references
func FindNode(position [3]int, potentialNodes [][][]*Node) *Node {
    return potentialNodes[position[0]][position[1]][position[2]]
}

func LoadState(name string, kb keyboard.Keyboard) *Network {
    fmt.Println(fmt.Sprintf("Loading state \"%v\"...", name))
    datafile, err := os.Open(fmt.Sprintf("./state/%v_state.json", name))
    if err != nil {
        fmt.Println(err)
    }
    decoder := json.NewDecoder(datafile)
    importedNet := &DisplayNetwork{}
    decoder.Decode(&importedNet)
    datafile.Close()

    net := &Network{
        Nodes: [][][]*Node{},
        Dimensions: importedNet.Dimensions,
    }
    // set nodes
    // this looks good
    importedNet.ForEachINode(func(importedNode *DisplayNode, pos [3]int) {
        newNode := &Node{
            Value: importedNode.Value,
            Position: importedNode.Position,
            IncomingConnections: []*Connection{},
        }
        net.Nodes[pos[0]][pos[1]][pos[2]] = newNode
    });
    // set connections
    // this part is super inefficient
    // still should optimize
    importedNet.ForEachINode(func(importedNode *DisplayNode, pos [3]int) {
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
    })

    // set sensors
    // this is also inefficient
    for _, importedSensor := range importedNet.Sensors {
        nodes := []*Node{}
        for _, nodePos := range importedSensor.Nodes {
            nodes = append(nodes, FindNode(nodePos, net.Nodes))
        }
        newSensor := &Sensor{
            Nodes: nodes,
            Excitatory: importedSensor.Excitatory,
            Trigger: importedSensor.Trigger,
            Stimulated: false,
            Name: importedSensor.Name,
        }
        net.Sensors = append(net.Sensors, newSensor)
        if kb != nil {
            kb.Bind(func() {
                newSensor.Stimulated = !newSensor.Stimulated
            }, importedSensor.Trigger)
        }
    }
    return net
}

func (net Network) SaveState(name string) {
    fmt.Println(fmt.Sprintf("Saving state \"%v\"...", name))
    os.Mkdir("state", 755)
    dispNet := DisplayNetwork{
        Nodes: [][][]*DisplayNode{},
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
            Name: sensor.Name,
        })
    }
    for i := 0; i < net.Dimensions[0]; i++ {
        iDim := [][]*DisplayNode{}
        for j := 0; j < net.Dimensions[1]; j++ {
            jDim := []*DisplayNode{}
            for k := 0; k < net.Dimensions[2]; k++ {
                node := net.Nodes[i][j][k]
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
                jDim = append(jDim, dispNode)
            }
            iDim = append(iDim, jDim)
        }
        dispNet.Nodes = append(dispNet.Nodes, iDim)
    }
    f, _ := os.Create(fmt.Sprintf("./state/%v_state.json", name))
    f.WriteString(dispNet.String())
    f.Close()
}

func (impNet DisplayNetwork) ForEachINode(handler func(*DisplayNode, [3]int)) {
    for i := range impNet.Nodes {
        for j := range impNet.Nodes[i] {
            for k := range impNet.Nodes[i][j] {
                handler(impNet.Nodes[i][j][k], [3]int{i, j, k})
            }
        }
    }
}