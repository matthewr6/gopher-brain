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

type DisplayNetwork struct {
    Nodes []*DisplayNode `json:"nodes"`
    Connections []*DisplayConnection `json:"connections"`
}

type DisplayNode struct {
    Value float64                             `json:"value"`
    Position [3]int                           `json:"position"`
}

type DisplayConnection struct {
    Strength float64      `json:"strength"`
    To [3]int             `json:"to"`
    From [3]int           `json:"from"`
    HoldingVal float64    `json:"holdingVal"`
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
    // manipulate network
    net := &Network{
        Nodes: []*Node{},
    }
    // set nodes
    for _, importedNode := range importedNet.Nodes {
        newNode := &Node{
            Value: importedNode.Value,
            Position: importedNode.Position,
            OutgoingConnections: []*Connection{},
            IncomingConnections: []*Connection{},
        }
        net.Nodes = append(net.Nodes, newNode)
    }
    // set connections
    for _, node := range net.Nodes {
        for _, potCon := range importedNet.Connections {
            if potCon.From == node.Position {
                node.OutgoingConnections = append(node.OutgoingConnections, &Connection{
                    Strength: potCon.Strength,
                    To: FindNode(potCon.To, net.Nodes),
                    HoldingVal: potCon.HoldingVal,
                })
            }
            if potCon.To == node.Position {
                node.IncomingConnections = append(node.IncomingConnections, &Connection{
                    Strength: potCon.Strength,
                    To: node,
                    HoldingVal: potCon.HoldingVal,
                })
            }
        }
    }
    return net
}

func (net Network) SaveState(name string) {
    fmt.Println("saving")
    dispNet := DisplayNetwork{
        Nodes: []*DisplayNode{},
        Connections: []*DisplayConnection{},
    }
    for _, node := range net.Nodes {
        dispNode := &DisplayNode{
            Value: node.Value,
            Position: node.Position,
        }
        for _, conn := range node.OutgoingConnections {
            to := [3]int{
                conn.To.Position[0],
                conn.To.Position[1],
                conn.To.Position[2],   
            }
            from := [3]int{
                node.Position[0],
                node.Position[1],
                node.Position[2],
            }
            dispConn := &DisplayConnection{
                Strength: conn.Strength,
                HoldingVal: conn.HoldingVal,
                To: to,
                From: from,
            }
            dispNet.Connections = append(dispNet.Connections, dispConn)
        }
        dispNet.Nodes = append(dispNet.Nodes, dispNode)
    }
    f, _ := os.Create(fmt.Sprintf("./%v_state.json", name))
    f.WriteString(dispNet.String())
    f.Close()
}