package brain

import (
    "fmt"
    "math"
    "bufio"
    "strconv"
    "strings"
)

func StrsToInts(strings []string) []int {
    retval := []int{}
    for _, i := range strings {
        j, _ := strconv.Atoi(i)
        retval = append(retval, j)
    }
    return retval
}

func IntDist(p1 [3]int, p2 [3]int) float64 {
    sum := 0
    for i := 0; i < 3; i++ {
        sum += (p1[i]-p2[i])*(p1[i]-p2[i])
    }
    return math.Sqrt(float64(sum))
}

func FloatDist(p1 [3]float64, p2 [3]float64) float64 {
    sum := 0.0
    for i := 0; i < 3; i++ {
        sum += (p1[i]-p2[i])*(p1[i]-p2[i])
    }
    return math.Sqrt(float64(sum))
}

func Prompt(message string, reader *bufio.Reader) string {
    fmt.Print(message)
    text, _ := reader.ReadString('\n')
    return strings.TrimSpace(text)
}

func (net Network) FindNode(position [3]int) *Node {
    return net.Nodes[position[0]][position[1]][position[2]]
}

func (net Network) FindLeftHemisphereNode(position [3]int) *Node {
    return net.LeftHemisphere[position[0]][position[1]][position[2]]
}

func (net Network) FindRightHemisphereNode(position [3]int) *Node {
    return net.RightHemisphere[position[0]][position[1]][position[2]]
}

func (net *Network) ForEachNode(handler func(*Node, [3]int)) {
    for i := range net.Nodes {
        for j := range net.Nodes[i] {
            for k := range net.Nodes[i][j] {
                handler(net.Nodes[i][j][k], [3]int{i, j, k})
            }
        }
    }
}


// returns (totalConnections, avgStrength, noInputs, noOutputs, isolatedNodes)
func (net *Network) CountConnections() (int, float64, int, int, int) {
    total := 0
    avgStrength := 0.0
    avgCounter := 0
    noInputs := 0
    noOutputs := 0
    isolated := 0
    net.ForEachNode(func(n *Node, pos [3]int) {
        total += len(n.OutgoingConnection.To)
        for _, info := range n.OutgoingConnection.To {
            avgStrength += info.Strength
            avgCounter += 1
        }
        if len(n.IncomingConnections) == 0 {
            noInputs += 1
        }
        if len(n.OutgoingConnection.To) == 0 {
            noOutputs += 1
        }
        if (len(n.IncomingConnections) + len(n.OutgoingConnection.To)) == 0 {
            isolated += 1
        }
    })
    return total, avgStrength/float64(avgCounter), noInputs, noOutputs, isolated
}

func (net *Network) ForEachRightHemisphereNode(handler func(*Node, [3]int)) {
    for i := range net.RightHemisphere {
        for j := range net.RightHemisphere[i] {
            for k := range net.RightHemisphere[i][j] {
                handler(net.RightHemisphere[i][j][k], [3]int{i, j, k})
            }
        }
    }
}

func (net *Network) ForEachLeftHemisphereNode(handler func(*Node, [3]int)) {
    for i := range net.LeftHemisphere {
        for j := range net.LeftHemisphere[i] {
            for k := range net.LeftHemisphere[i][j] {
                handler(net.LeftHemisphere[i][j][k], [3]int{i, j, k})
            }
        }
    }
}