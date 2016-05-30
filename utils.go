package main

import (
    "strconv"
    "bufio"
    "strings"
    "fmt"
)

func StrsToInts(strings []string) []int {
    retval := []int{}
    for _, i := range strings {
        j, _ := strconv.Atoi(i)
        retval = append(retval, j)
    }
    return retval
}

func Prompt(message string, reader *bufio.Reader) string {
    fmt.Print(message)
    text, _ := reader.ReadString('\n')
    return strings.TrimSpace(text)
}

func FindNode(position [3]int, potentialNodes [][][]*Node) *Node {
    return potentialNodes[position[0]][position[1]][position[2]]
}