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
        node := FindNode(importedNode.Position, net.Nodes)
        nodeToConnect := FindNode(importedNode.OutgoingConnection.To, net.Nodes)
        newConn := &Connection{
            HoldingVal: importedNode.OutgoingConnection.HoldingVal,
            Terminals: importedNode.OutgoingConnection.Terminals,
            Excitatory: importedNode.OutgoingConnection.Excitatory,
            To: nodeToConnect,
        }
        node.OutgoingConnection = newConn
        nodeToConnect.IncomingConnections = append(nodeToConnect.IncomingConnections, newConn)
    }
    return net
}

func (net Network) SaveState(name string) {
    fmt.Println("saving")
    dispNet := DisplayNetwork{
        Nodes: []*DisplayNode{},
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