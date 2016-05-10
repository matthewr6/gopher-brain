package main

import (
    "fmt"
    "time"
    // "os"
    // "strconv"
    // "reflect"
    "math/rand"
)

/*
    & is "address of"
    * is "value at address"
*/

func main() {
    start := time.Now()
    rand.Seed(time.Now().UTC().UnixNano())

    // [width, depth, height]
    NETWORK_SIZE := [3]int{25, 25, 25}
    myNet := MakeNetwork(NETWORK_SIZE, false)
    myNet.Connect()
    myNet.CreateSensor(1, 20, "x", [3]int{1, 5, 5})

    // myNet := LoadState("test")
    // myNet.Stimulate([]Stimulus{
    //     Stimulus{
    //         Position: [3]int{1,1,1},
    //     },
    //     Stimulus{
    //         Position: [3]int{1,1,2},
    //     },
    //     Stimulus{
    //         Position: [3]int{1,2,1},
    //     },
    //     Stimulus{
    //         Position: [3]int{2,1,1},
    //     },
    //     Stimulus{
    //         Position: [3]int{1,2,2},
    //     },
    //     Stimulus{
    //         Position: [3]int{2,1,2},
    //     },
    //     Stimulus{
    //         Position: [3]int{2,2,1},
    //     },
    //     Stimulus{
    //         Position: [3]int{2,2,2},
    //     },
    // })
    // myNet.RandomizeValues()

    // frames, err := strconv.Atoi(os.Args[1])
    // if err != nil {
    //     fmt.Println(err)
    //     return
    // }
    // myNet.GenerateAnim(frames)
    
    // myNet.SaveState("test")
    
    // myNet.SaveState("test")
    // loadedNet := LoadState("test")
    // fmt.Println(reflect.DeepEqual(loadedNet, myNet))

    elapsed := time.Since(start)
    fmt.Printf("Took %s\n", elapsed)
}