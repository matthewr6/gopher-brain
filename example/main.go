package main

// import "github.com/firedrake969/gopher-brain"
import (
    ".."

    "os"
    "fmt"
    "bufio"
    "strings"
    "strconv"
)

func main() {
    reader := bufio.NewReader(os.Stdin)

    myNet := brain.Brain([3]int{12, 25, 25}, []brain.SensorConstructor{
        brain.SensorConstructor{
            Name:"eye",
            R: 1,
            Count: 9,
            Plane: "y",
            Center: [3]int{8, 0, 12},
            OutputCount: 2,
            InputFunc: func(nodes []*brain.Node, influences map[string]*brain.Output) {
                for _, node := range nodes {
                    node.Value = 1
                }
            },
        },
    })

    framesInput := brain.Prompt("Enter number of frames, or leave blank to run until manually stopped:  ", reader)
    frames, err := strconv.Atoi(framesInput)
    if err != nil {
        frames = 0
    }

    directory := brain.Prompt("Enter directory to save frames and state to:  ", reader)
    if directory == "" {
        directory = "."
    }
    if directory[len(directory)-1] == '/' {
        directory = directory[0:len(directory)-1]
    }
    
    myNet.GenerateAnim(frames)

    fmt.Print("\nSave state?  Enter a name if you wish to save the state:  ")
    fileName, _ := reader.ReadString('\n')
    fileName = strings.TrimSpace(fileName)
    if fileName != "" {
        myNet.SaveState(fileName)
    }
}