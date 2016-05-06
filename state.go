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
    From [3]int           `json:"from"`
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
    for _, importedNode := range importedNet.Nodes {
        // newConn := &Connection{
        //     HoldingVal: importedNode.OutgoingConnection.HoldingVal,
        //     Terminals: importedNode.OutgoingConnection.Terminals,
        //     Excitatory: importedNode.OutgoingConnection.Excitatory,
        // }
        newNode := &Node{
            Value: importedNode.Value,
            Position: importedNode.Position,
            // OutgoingConnection: newConn,
            IncomingConnections: []*Connection{},
        }
        net.Nodes = append(net.Nodes, newNode)
    }
    // todo
    // set connections
    // this part is super inefficient
    // still should optimize

    // two options:
    //     - iterate through all nodes and through all nodes inside that
    //         - bleh, NODE_AMOUNT * NODE_AMOUNT
    //         - may be harder than the latter
    //     - iterate through all nodes and iterate through all possible connections - what the old LoadState func did
    //         - bleh, NODE_AMOUNT * CONNECTION_AMOUNT
    // one outgoing connection per node
    // therefore NODE_AMOUNT = CONNECTION_AMOUNT
    // both approaches should then be the same speed
    // must determine speed
    // maybe have IDs to tag nodes?
    //     would require more complexity on the SaveState though
    for _,  importedNode := range importedNet.Nodes {
        fromNode := FindNode(importedNode.Position, net.Nodes)
        toNode := FindNode(importedNode.OutgoingConnection.To, net.Nodes)
        newConn := &Connection{
            HoldingVal: importedNode.OutgoingConnection.HoldingVal,
            Terminals: importedNode.OutgoingConnection.Terminals,
            Excitatory: importedNode.OutgoingConnection.Excitatory,
            To: toNode,
        }
        fromNode.OutgoingConnection = newConn
        toNode.IncomingConnections = append(fromNode.IncomingConnections, newConn)
    }
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
            From: node.Position,
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