package brain

import (
    "strconv"
    "bufio"
    "strings"
    "fmt"
    "math"
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