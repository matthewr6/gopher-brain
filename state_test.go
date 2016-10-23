package main

import (
    "fmt"
    "testing"
)

func TestState(t *testing.T) {
    NETWORK_SIZE := [3]int{25, 25, 25}
    testingNet := MakeNetwork(NETWORK_SIZE, false)
    testingNet.Connect()
    testingNet.Mirror()
    testingNet.ConnectHemispheres()

    testingNet.CreateSensor("eye", 1, 9, "y", [3]int{8, 0, 12}, 2, func(nodes []*Node, influences map[string]*Output) {
        for _, node := range nodes {
            node.Value = 1
        }
    })

    // this is to make sure adding/removing connections works
    for i := 0; i < 100; i++ {
        testingNet.Cycle()
    }

    testingNet.SaveState("test")
    loadedNet := LoadState("test")
    fmt.Println("Finished loading state.")
    same := Test(testingNet, loadedNet)
    if !same {
        t.Error("Loaded state did not match original state.")
    }
}

func BenchmarkBuildNet(b *testing.B) {
    NETWORK_SIZE := [3]int{25, 25, 25}
    for i := 0; i < b.N; i++ {
        testingNet := MakeNetwork(NETWORK_SIZE, false)
        testingNet.Connect()
        testingNet.Mirror()
        testingNet.ConnectHemispheres()
    }
}

func BenchmarkCycleNet(b *testing.B) {
    NETWORK_SIZE := [3]int{25, 25, 25}
    testingNet := MakeNetwork(NETWORK_SIZE, false)
    testingNet.Connect()
    testingNet.Mirror()
    testingNet.ConnectHemispheres()
    for i := 0; i < b.N; i++ {
        testingNet.Cycle()
    }
}