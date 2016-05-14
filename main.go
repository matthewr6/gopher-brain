package main

import (
    "fmt"
    "time"
    // "os"
    // "strconv"
    // "reflect"
    "math/rand"

    "github.com/jteeuwen/keyboard/termbox"
    term "github.com/nsf/termbox-go"
)

/*
    & is "address of"
    * is "value at address"
*/

var running = true

// todo - clean up goroutines and whatnot
func main() {
    start := time.Now()
    rand.Seed(time.Now().UTC().UnixNano())

    // [width, depth, height]
    NETWORK_SIZE := [3]int{25, 25, 25}
    myNet := MakeNetwork(NETWORK_SIZE, false)
    myNet.Connect()
    // myNet.CreateSensor(3, 25, "", [3]int{1, 1, 1}, true)
    // myNet.CreateSensor(2, 25, "y", [3]int{15, 1, 15}, true)

    // this is the keyboard sensing stuff
    term.Init()
    kb := termbox.New()
    kb.Bind(func() {
        running = false
    }, "space")
    go KeyboardPoll(kb)

    myNet.CreateSensor(3, 50, "", [3]int{25, 1, 1}, true, "a", kb)
    myNet.CreateSensor(1, 50, "", [3]int{1, 1, 1}, true, "s", kb)
    myNet.AnimateUntilDone(100)

    // myNet := LoadState("test")
    // todo - just make this an array of 3-length int arrays
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
    
    // myNet.SaveState("environ")
    
    // myNet.SaveState("test")
    // loadedNet := LoadState("test")
    // fmt.Println(reflect.DeepEqual(loadedNet, myNet))

    elapsed := time.Since(start)
    term.Close()
    fmt.Printf("Took %s\n", elapsed)
}