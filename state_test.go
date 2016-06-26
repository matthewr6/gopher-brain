package main

import (
    "fmt"
    "testing"
)

func TestState(t *testing.T) {
    NETWORK_SIZE := [3]int{2, 2, 2}
    testingNet := MakeNetwork(NETWORK_SIZE, false)
    testingNet.Connect()
    testingNet.Mirror()

    // testingNet.CreateSensor("aa", 1, 50, "", [3]int{24, 0, 0}, true, "a", func(nodes []*Node, stimulated bool) {
    //     for _, node := range nodes {
    //         if stimulated {
    //             node.Value = 1
    //         }
    //     }
    // })
    // testingNet.CreateSensor("bb", 1, 50, "", [3]int{0, 0, 0}, true, "b", func(nodes []*Node, stimulated bool) {
    //     for _, node := range nodes {
    //         if stimulated {
    //             node.Value = 1
    //         }
    //     }
    // })
    // testingNet.CreateOutput("output", 1, 50,"", [3]int{12, 1, 1}, func(nodes []*Node) float64 {
    //     var sum float64
    //     for _, node := range nodes {
    //         if node.OutgoingConnection.To[node].Excitatory {
    //             sum += float64(node.Value) * node.OutgoingConnection.To[node].Strength
    //         } else {
    //             sum -= float64(node.Value) * node.OutgoingConnection.To[node].Strength
    //         }
    //     }
    //     return sum
    // })

    testingNet.SaveState("test")
    loadedNet := LoadState("test")
    fmt.Println("Finished loading state.")
    same := Test(testingNet, loadedNet)
    if !same {
        t.Error("Loaded state did not match original state.")
    }
}