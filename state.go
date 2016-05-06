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
    Nodes []*DisplayNode `json:"nodes"`
    // Connections []*DisplayConnection `json:"connections"`
}

type DisplayNode struct {
    Value int                               `json:"value"`
    Position [3]int                         `json:"position"`
    OutgoingConnection *DisplayConnection   `json:"axon"`
}

type DisplayConnection struct {
    To [3]int             `json:"to"`
    HoldingVal int        `json:"holdingVal"`
    Terminals int         `json:"terminals"`
    Excitatory bool       `json:"excitatory"`
}

func (d DisplayNetwork) String() string {
    jsonRep, _ := json.MarshalIndent(d, "", "    ")
    return string(jsonRep)
}

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
    datafile, err := os.Open(fmt.Sprintf("./%v_state.json", name))
    if err != nil {
        fmt.Println(err)
    }
    decoder := json.NewDecoder(datafile)
    importedNet := &DisplayNetwork{}
    decoder.Decode(&importedNet)
    datafile.Close()
    net := &Network{
        Nodes: []*Node{},
    }
    // set nodes
    for _, importedNode := range importedNet.Nodes {
        newConn := &Connection{
            HoldingVal: importedNode.OutgoingConnection.HoldingVal,
            Terminals: importedNode.OutgoingConnection.Terminals,
            Excitatory: importedNode.OutgoingConnection.Excitatory,
        }
        newNode := &Node{
            Value: importedNode.Value,
            Position: importedNode.Position,
            OutgoingConnection: newConn,
            IncomingConnections: []*Connection{},
        }
        net.Nodes = append(net.Nodes, newNode)
    }
    // todo
    // set connections
    // this part is super inefficient
    // for _, node := range net.Nodes {
    //     for _, potCon := range importedNet.Connections {
    //         if potCon.From == node.Position {
    //             node.OutgoingConnection = &Connection{
    //                 To: FindNode(potCon.To, net.Nodes),
    //                 HoldingVal: potCon.HoldingVal,
    //             }
    //         }
    //         if potCon.To == node.Position {
    //             node.IncomingConnections = append(node.IncomingConnections, &Connection{
    //                 To: node,
    //                 HoldingVal: potCon.HoldingVal,
    //             })
    //         }
    //     }
    // }
    return net
}

func (net Network) SaveState(name string) {
    fmt.Println("saving")
    dispNet := DisplayNetwork{
        Nodes: []*DisplayNode{},
        // Connections: []*DisplayConnection{},
    }
    for _, node := range net.Nodes {
        dispConn := &DisplayConnection{
            To: node.OutgoingConnection.To.Position,
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
    f, _ := os.Create(fmt.Sprintf("./%v_state.json", name))
    f.WriteString(dispNet.String())
    f.Close()
}