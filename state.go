package main

import (
    "os"
    "fmt"
    "strings"
    "reflect"
    "encoding/json"
)

/*
    & is "address of"
    * is "value at address"
*/

// still have to reconcile the updated network.go stuff

type DisplayNetwork struct {
    Nodes [][][]*DisplayNode            `json:"nodes"`
    Dimensions [3]int                   `json:"dimensions"`
    Sensors map[string]*DisplaySensor   `json:"sensors"`
    Outputs map[string]*DisplayOutput   `json:"outputs"`
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
    Center [3]int           `json:"center"`
}

// // do I want to save sensors/outputs?
// type Output struct {
//     Nodes map[*Node]*ConnInfo               `json:"nodes"`
//     Name string                             `json:"name"`
//     Value float64                           `json:"value"`
//     Out func(map[*Node]*ConnInfo) float64   `json:"-"`
// }

type DisplaySensor struct {
    Nodes [][3]int        `json:"nodes"`
    Influences []string   `json:"influences"`
    Center [3]int         `json:"center"`
    Name string           `json:"name"`
}

type DisplayOutput struct {
    Nodes map[string]*ConnInfo    `json:"nodes"` // why pointers?  oh well it works so yeah
    Name string                   `json:"name"`
    Value float64                 `json:"value"`
}

func (d DisplayNetwork) String() string {
    jsonRep, _ := json.MarshalIndent(d, "", "    ")
    return string(jsonRep)
}

func (d DisplayOutput) String() string {
    jsonRep, _ := json.MarshalIndent(d, "", "    ")
    return string(jsonRep)
}

func (d DisplaySensor) String() string {
    jsonRep, _ := json.MarshalIndent(d, "", "    ")
    return string(jsonRep)
}

func LoadState(name string) *Network {
    fmt.Println(fmt.Sprintf("Loading state \"%v\"...", name))
    datafile, err := os.Open(fmt.Sprintf("%v/state/%v_state.json", directory, name))
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
        Sensors: make(map[string]*Sensor),
        Outputs: make(map[string]*Output),
    }
    for i := 0; i < (net.Dimensions[0]*2); i++ {
        iDim := [][]*Node{}
        for j := 0; j < net.Dimensions[1]; j++ {
            jDim := []*Node{}
            for k := 0; k < net.Dimensions[2]; k++ {
                newNode := &Node{
                    Value: importedNet.Nodes[i][j][k].Value,
                    Position: importedNet.Nodes[i][j][k].Position,
                    IncomingConnections: make(map[*Node]*Connection),
                    Id: fmt.Sprintf("%v|%v|%v", i, j, k),
                }
                jDim = append(jDim, newNode)
            }
            iDim = append(iDim, jDim)
        }
        net.Nodes = append(net.Nodes, iDim)
    }
    importedNet.ForEachINode(func(importedNode *DisplayNode, pos [3]int) {
        newConn := &Connection{
            HoldingVal: importedNode.OutgoingConnection.HoldingVal,
            Center: importedNode.OutgoingConnection.Center,
        }
        node := net.FindNode(importedNode.Position)
        toNodes := make(map[*Node]*ConnInfo)
        for id, connInfo := range importedNode.OutgoingConnection.To {
            posSlice := StrsToInts(strings.Split(id, "|"))
            nodeToConnect := net.FindNode([3]int{posSlice[0], posSlice[1], posSlice[2]})
            toNodes[nodeToConnect] = connInfo
            nodeToConnect.IncomingConnections[node] = newConn
        }
        newConn.To = toNodes
        node.OutgoingConnection = newConn
    })

    // WORKING

    // first, load outputs
    for _, importedOutput := range importedNet.Outputs {
        newOutput := &Output{
            Name: importedOutput.Name,
            Nodes: make(map[*Node]*ConnInfo),
            Value: importedOutput.Value,
            Out: func(nodes map[*Node]*ConnInfo) float64 {
                var sum float64
                for node, connInfo := range nodes {
                    if connInfo.Excitatory {
                        sum += float64(node.Value) * connInfo.Strength
                    } else {
                        sum -= float64(node.Value) * connInfo.Strength
                    }
                }
                return sum
            },
        }
        for id, info := range importedOutput.Nodes {
            pos := StrsToInts(strings.Split(id, "|"))
            node := net.FindNode([3]int{pos[0], pos[1], pos[2]})
            newConnInfo := &ConnInfo{
                Excitatory: info.Excitatory,
                Strength: info.Strength,
            }
            newOutput.Nodes[node] = newConnInfo
        }
        net.Outputs[importedOutput.Name] = newOutput
    }

    // then load sensors
    for _, importedSensor := range importedNet.Sensors {
        newSensor := &Sensor{
            Name: importedSensor.Name,
            Nodes: []*Node{},
            Influences: make(map[string]*Output),
            Center: importedSensor.Center,
        }
        for _, influenceName := range importedSensor.Influences {
            newSensor.Influences[influenceName] = net.Outputs[influenceName]
        }
        for _, nodePos := range importedSensor.Nodes {
            newSensor.Nodes = append(newSensor.Nodes, net.FindNode(nodePos))
        }
        net.Sensors[importedSensor.Name] = newSensor
    }

    // END WORKING

    return net
}

func (net Network) SaveState(name string) {
    fmt.Println(fmt.Sprintf("Saving state \"%v\"...", name))
    os.Mkdir("state", 755)
    dispNet := DisplayNetwork{
        Nodes: [][][]*DisplayNode{},
        Sensors: make(map[string]*DisplaySensor),
        Outputs: make(map[string]*DisplayOutput),
        Dimensions: net.Dimensions,
    }
    for i := 0; i < (net.Dimensions[0]*2); i++ {
        iDim := [][]*DisplayNode{}
        for j := 0; j < net.Dimensions[1]; j++ {
            jDim := []*DisplayNode{}
            for k := 0; k < net.Dimensions[2]; k++ {
                node := net.Nodes[i][j][k]
                toNodes := make(map[string]*ConnInfo)
                for connNode, connInfo := range node.OutgoingConnection.To {
                    toNodes[connNode.Id] = connInfo
                }
                dispConn := &DisplayConnection{
                    To: toNodes,
                    HoldingVal: node.OutgoingConnection.HoldingVal,
                    Center: node.OutgoingConnection.Center,
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

    // do we even need to reference the outputs when saving a state?   based on the way the outputs are created, it's possible that we only need to reference the outputs when building the sensors when we load the sensors into the working net
    for _, output := range net.Outputs {
        nodeMap := make(map[string]*ConnInfo)
        for node, connInfo := range output.Nodes {
            nodeMap[node.Id] = connInfo
        }
        dispNet.Outputs[output.Name] = &DisplayOutput{
            Nodes: nodeMap,
            Name: output.Name,
        }
    }
    for _, sensor := range net.Sensors {
        positions := [][3]int{}
        for _, sensoryNode := range sensor.Nodes {
            positions = append(positions, sensoryNode.Position)
        }
        influenceNames := []string{}
        for _, influence := range sensor.Influences {
            influenceNames = append(influenceNames, influence.Name)
        }
        dispNet.Sensors[sensor.Name] = &DisplaySensor{
            Nodes: positions,
            Name: sensor.Name,
            Center: sensor.Center,
            Influences: influenceNames,
        }
    }

    f, _ := os.Create(fmt.Sprintf("%v/state/%v_state.json", directory, name))
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
    // don't need to compare left/right hemis
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
                    // fmt.Println("here")
                    return false
                }


                // compare outgoing connections
                // - compare immediate connection data
                if (oNode.OutgoingConnection.HoldingVal != lNode.OutgoingConnection.HoldingVal || 
                    oNode.OutgoingConnection.Center != lNode.OutgoingConnection.Center) {
                    return false
                }
                // - compare To data
                //   - https://groups.google.com/forum/#!topic/golang-nuts/UWKAOXyMwJM
                oOutgoingTo := make(map[string]*ConnInfo)
                lOutgoingTo := make(map[string]*ConnInfo)
                for node, info := range oNode.OutgoingConnection.To {
                    oOutgoingTo[node.Id] = info
                }
                for node, info := range lNode.OutgoingConnection.To {
                    lOutgoingTo[node.Id] = info
                }
                if !reflect.DeepEqual(oOutgoingTo, lOutgoingTo) {
                    return false
                }

                // compare incoming connections
                oIncoming := make(map[string]*ConnInfo)
                lIncoming := make(map[string]*ConnInfo)
                for from, conn := range oNode.IncomingConnections {
                    oIncoming[from.Id] = conn.To[oNode]
                }
                for from, conn := range lNode.IncomingConnections {
                    lIncoming[from.Id] = conn.To[lNode]
                }
                // - compare number of incoming connections
                if (len(oNode.IncomingConnections) != len(lNode.IncomingConnections)) {
                    return false
                }
                // - compare incoming connection info
                if !reflect.DeepEqual(oIncoming, lIncoming) {
                    for i := range oIncoming {
                        fmt.Println(oIncoming[i], lIncoming[i])
                    }
                    return false
                }
            }
        }
    }

    // to compare
    // sensors
    //     - influence names
    //     - node positions
    //     - names (compare *based on* names)
    //     - centers
    // outputs
    //     - names (compare based on names)
    //     - nodes (compare node positions and conninfo stuff)
    //     - values
    for sensorName, oSensor := range orig.Sensors {
        lSensor, exists := loaded.Sensors[sensorName]
        if !exists {
            return false
        }
        if ((oSensor.Center != lSensor.Center) ||
            (len(oSensor.Nodes) != len(lSensor.Nodes)) ||
            (len(oSensor.Influences) != len(lSensor.Influences))) {
            return false
        }
        for idx, oNode := range oSensor.Nodes {
            lNode := lSensor.Nodes[idx]
            if oNode.Id != lNode.Id {
                return false
            }
        }
    }
    for outputName, oOutput := range orig.Outputs {
        lOutput, exists := loaded.Outputs[outputName]
        if !exists {
            return false
        }
        if ((len(oOutput.Nodes) != len(lOutput.Nodes)) ||
            (oOutput.Value != lOutput.Value)) {
            return false
        }
        oOutputNodes := make(map[string]float64)
        lOutputNodes := make(map[string]float64)
        for node, info := range oOutput.Nodes {
            oOutputNodes[node.Id] = info.Strength
        }
        for node, info := range lOutput.Nodes {
            lOutputNodes[node.Id] = info.Strength
        }
        if !reflect.DeepEqual(oOutputNodes, lOutputNodes) {
            return false
        }
    }
    return true
}