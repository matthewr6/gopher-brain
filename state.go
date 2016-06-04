package main

import (
    "os"
    "fmt"
    "strings"
    "reflect"
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
    Outputs []*DisplayOutput   `json:"outputs"`
    // Connections []*DisplayConnection `json:"connections"`
}

type DisplayNode struct {
    Value int                               `json:"value"`
    Position [3]int                         `json:"position"`
    OutgoingConnection *DisplayConnection   `json:"axon"`
    Id string                               `json:"id"`
}

type DisplayConnection struct {
    To map[string]*ConnInfo `json:"to"` // shouldn't need a display struct for this, also idk why it's *ConnInfo instead of ConnInfo but hey it works
    HoldingVal int          `json:"holdingVal"`
}

type DisplaySensor struct {
    Nodes [][3]int    `json:"nodes"`
    Excitatory bool   `json:"excitatory"`
    Trigger string    `json:"trigger"`
    Name string       `json:"name"`
}

type DisplayOutput struct {
    Nodes map[string]*ConnInfo    `json:"nodes"` // why pointers?  oh well it works so yeah
    Name string                   `json:"name"`
}

func (d DisplayNetwork) String() string {
    jsonRep, _ := json.MarshalIndent(d, "", "    ")
    return string(jsonRep)
}

// func LoadState(name string) *Network {
//     fmt.Println(fmt.Sprintf("Loading state \"%v\"...", name))
//     datafile, err := os.Open(fmt.Sprintf("./state/%v_state.json", name))
//     if err != nil {
//         fmt.Println(err)
//     }
//     decoder := json.NewDecoder(datafile)
//     importedNet := &DisplayNetwork{}
//     decoder.Decode(&importedNet)
//     datafile.Close()

//     net := &Network{
//         Nodes: [][][]*Node{},
//         Dimensions: importedNet.Dimensions,
//     }
//     // set nodes
//     // this looks good
//     for i := 0; i < net.Dimensions[0]; i++ {
//         iDim := [][]*Node{}
//         for j := 0; j < net.Dimensions[1]; j++ {
//             jDim := []*Node{}
//             for k := 0; k < net.Dimensions[2]; k++ {
//                 newNode := &Node{
//                     Value: importedNet.Nodes[i][j][k].Value,
//                     Position: importedNet.Nodes[i][j][k].Position,
//                     IncomingConnections: []*Connection{},
//                     Id: fmt.Sprintf("%v|%v|%v", i, j, k),
//                 }
//                 jDim = append(jDim, newNode)
//             }
//             iDim = append(iDim, jDim)
//         }
//         net.Nodes = append(net.Nodes, iDim)
//     }
//     // set connections
//     // this part is super inefficient
//     // still should optimize
//     importedNet.ForEachINode(func(importedNode *DisplayNode, pos [3]int) {
//         newConn := &Connection{
//             HoldingVal: importedNode.OutgoingConnection.HoldingVal,
//         }
//         node := FindNode(importedNode.Position, net.Nodes)
//         toNodes := make(map[*Node]*ConnInfo)
//         for id, connInfo := range importedNode.OutgoingConnection.To {
//             posSlice := StrsToInts(strings.Split(id, "|"))
//             nodeToConnect := FindNode([3]int{posSlice[0], posSlice[1], posSlice[2]}, net.Nodes)
//             toNodes[nodeToConnect] = connInfo
//             nodeToConnect.IncomingConnections = append(nodeToConnect.IncomingConnections, newConn)
//         }
//         newConn.To = toNodes
//         node.OutgoingConnection = newConn
//     })

//     // set sensors
//     // this is also inefficient
//     for _, importedSensor := range importedNet.Sensors {
//         nodes := []*Node{}
//         for _, nodePos := range importedSensor.Nodes {
//             nodes = append(nodes, FindNode(nodePos, net.Nodes))
//         }
//         newSensor := &Sensor{
//             Nodes: nodes,
//             Excitatory: importedSensor.Excitatory,
//             Trigger: importedSensor.Trigger,
//             Stimulated: false,
//             Name: importedSensor.Name,
//             In: func(nodes []*Node, stimulated bool) {
//                 // for simplicity - just continuously stimulate every node
//                 for _, node := range nodes {
//                     if stimulated {
//                         node.Value = 1
//                     }
//                     // let's try removing this for now, see what happens...
//                     // else {
//                     //     node.Value = 0
//                     // }
//                 }
//             },
//         }
//         net.Sensors = append(net.Sensors, newSensor)
//         // if kb != nil {
//         //     kb.Bind(func() {
//         //         newSensor.Stimulated = !newSensor.Stimulated
//         //     }, importedSensor.Trigger)
//         // }
//     }

//     for _, importedOutput := range importedNet.Outputs {
//         nodes := []*Node{}
//         for _, nodePos := range importedOutput.Nodes {
//             nodes = append(nodes, FindNode(nodePos, net.Nodes))
//         }
//         newOutput := &Output{
//             Nodes: nodes,
//             Name: importedOutput.Name,
//             Out: func(nodes []*Node) float64 {
//                 var sum float64
//                 for _, node := range nodes {
//                     if node.OutgoingConnection.To[node].Excitatory {
//                         sum += float64(node.Value) * node.OutgoingConnection.To[node].Strength
//                     } else {
//                         sum -= float64(node.Value) * node.OutgoingConnection.To[node].Strength
//                     }
//                 }
//                 return sum
//             },
//         }
//         net.Outputs = append(net.Outputs, newOutput)
//     }

//     return net
// }

func (net *Network) BindKeyboard(kb keyboard.Keyboard) {
    for _, sensor := range net.Sensors {
        s := sensor
        kb.Bind(func() {
            s.Stimulated = !s.Stimulated
        }, sensor.Trigger)
    }
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

    for _, output := range net.Outputs {
        // positions := [][3]int{}
        nodeMap := make(map[string]*ConnInfo)
        for node, connInfo := range output.Nodes {
            nodeMap[node.Id] = connInfo
        }
        // for _, outputNode := range output.Nodes {
        //     positions = append(positions, outputNode.Position)
        // }
        dispNet.Outputs = append(dispNet.Outputs, &DisplayOutput{
            Nodes: nodeMap,
            Name: output.Name,
        })
    }

    for i := 0; i < net.Dimensions[0]; i++ {
        iDim := [][]*DisplayNode{}
        for j := 0; j < net.Dimensions[1]; j++ {
            jDim := []*DisplayNode{}
            for k := 0; k < net.Dimensions[2]; k++ {
                node := net.Nodes[i][j][k]
                toNodes := make(map[string]*ConnInfo)
                for node, connInfo := range node.OutgoingConnection.To {
                    toNodes[node.Id] = connInfo
                }
                dispConn := &DisplayConnection{
                    To: toNodes,
                    HoldingVal: node.OutgoingConnection.HoldingVal,
                    // Excitatory: node.OutgoingConnection.Excitatory,
                    // Strength: node.OutgoingConnection.Strength,
                }

                dispNode := &DisplayNode{
                    Value: node.Value,
                    Position: node.Position,
                    OutgoingConnection: dispConn,
                    Id: node.Id,
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

func Test(orig, loaded *Network) bool {
    // compare dimensions
    if orig.Dimensions != loaded.Dimensions {
        return false
    }
    // compare nodes
    for i := range orig.Nodes {
        for j := range orig.Nodes[i] {
            for k := range orig.Nodes[i][j] {
                oNode := orig.Nodes[i][j][k]
                lNode := loaded.Nodes[i][j][k]
                // first compare value, position, id
                if ((oNode.Value != lNode.Value) ||
                    (oNode.Position != lNode.Position) ||
                    (oNode.Id != lNode.Id)) {
                    return false
                }
                // then compare the input/output connections
                oConns := []Connection{*oNode.OutgoingConnection}
                for _, iConn := range oNode.IncomingConnections {
                    oConns = append(oConns, *iConn)
                }

                lConns := []Connection{*lNode.OutgoingConnection}
                for _, lConn := range lNode.IncomingConnections {
                    lConns = append(lConns, *lConn)
                }

                if len(oConns) != len(lConns) {
                    return false
                }
                // this should work because slices and arrays are ordered
                for i := range oConns {
                    if oConns[i].HoldingVal != lConns[i].HoldingVal {
                        return false
                    }
                }
                // convert the maps to map[string]*ConnInfo to be comparable
                oConnInfo := map[string]*ConnInfo{}
                for _, conn := range oConns {
                    for node, connInfo := range conn.To {
                        oConnInfo[node.Id] = connInfo
                    }
                }
                lConnInfo := map[string]*ConnInfo{}
                for _, conn := range lConns {
                    for node, connInfo := range conn.To {
                        lConnInfo[node.Id] = connInfo
                    }
                }
                if !reflect.DeepEqual(oConnInfo, lConnInfo) {
                    return false
                }
            }
        }
    }
    // these still use reflect for now
    // compare sensors
    if !reflect.DeepEqual(orig.Sensors, loaded.Sensors) {
        return false
    }
    // compare outputs
    if !reflect.DeepEqual(orig.Outputs, loaded.Outputs) {
        return false
    }
    return true
}