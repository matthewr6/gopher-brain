package main

import (
    "fmt"
    // "os"
    "time"
    // "strconv"
    "reflect"
)

/*
    & is "address of"
    * is "value at address"
*/

func main() {
    start := time.Now()

    // [width, depth, height]
    NETWORK_SIZE := [3]int{2, 2, 2}
    myNet := MakeNetwork(NETWORK_SIZE, false)
    myNet.Connect()

    // myNet := LoadState("test")
    // myNet.Stimulate([]Stimulus{
    //     Stimulus{
    //         Position: [3]int{25,1,1},
    //         Strength: 5,
    //     },
    //     Stimulus{
    //         Position: [3]int{24,1,1},
    //         Strength: 5,
    //     },
    //     Stimulus{
    //         Position: [3]int{25,2,1},
    //         Strength: 5,
    //     },
    //     Stimulus{
    //         Position: [3]int{25,1,2},
    //         Strength: 5,
    //     },
    //     Stimulus{
    //         Position: [3]int{25,2,2},
    //         Strength: 5,
    //     },
    //     Stimulus{
    //         Position: [3]int{24,1,2},
    //         Strength: 5,
    //     },
    //     Stimulus{
    //         Position: [3]int{24,2,1},
    //         Strength: 5,
    //     },
    //     Stimulus{
    //         Position: [3]int{24,2,2},
    //         Strength: 5,
    //     },
    //     Stimulus{
    //         Position: [3]int{23, 1, 1},
    //         Strength: 5,
    //     },
    //     Stimulus{
    //         Position: [3]int{21, 1, 1},
    //         Strength: 5,
    //     },
    // })
    // myNet.RandomizeValues()

    // frames, err := strconv.Atoi(os.Args[1])
    // if err != nil {
    //     fmt.Println(err)
    //     return
    // }
    // myNet.GenerateAnim(frames)
    
    // myNet.SaveState("test2")
    myNet.SaveState("test")
    fmt.Println(reflect.DeepEqual(myNet, LoadState("test")))

    elapsed := time.Since(start)
    fmt.Printf("Took %s\n", elapsed)
}