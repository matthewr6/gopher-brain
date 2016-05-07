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
    NETWORK_SIZE := [3]int{2, 2, 1}
    myNet := MakeNetwork(NETWORK_SIZE, false)
    myNet.Connect()

    // myNet := LoadState("test")
    // myNet.Stimulate([]Stimulus{
    //     Stimulus{
    //         Position: [3]int{25,1,1},
    //     },
    //     Stimulus{
    //         Position: [3]int{24,1,1},
    //     },
    //     Stimulus{
    //         Position: [3]int{25,2,1},
    //     },
    //     Stimulus{
    //         Position: [3]int{25,1,2},
    //     },
    //     Stimulus{
    //         Position: [3]int{25,2,2},
    //     },
    //     Stimulus{
    //         Position: [3]int{24,1,2},
    //     },
    //     Stimulus{
    //         Position: [3]int{24,2,1},
    //     },
    //     Stimulus{
    //         Position: [3]int{24,2,2},
    //     },
    //     Stimulus{
    //         Position: [3]int{23, 1, 1},
    //     },
    //     Stimulus{
    //         Position: [3]int{21, 1, 1},
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
    // it could be that some data is missing when it's saved vs when it's created
    // could it be that incoming connections aren't getting set?
    // that's the only invisible thing, so possibly...
    myNet.SaveState("test")
    loadedNet := LoadState("test")
    fmt.Println(reflect.DeepEqual(loadedNet, myNet))

    elapsed := time.Since(start)
    fmt.Printf("Took %s\n", elapsed)
}