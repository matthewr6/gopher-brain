package main

// import ".."
import (
    "github.com/matthewr6/gopher-brain"

    "os"
    "fmt"
    "bufio"
    "strings"
    "strconv"
)

var directory = ""

func main() {
    reader := bufio.NewReader(os.Stdin)

    myNet := brain.Brain([3]int{12, 12, 12}, []brain.SensorConstructor{
        brain.SensorConstructor{
            Name:"eye",
            R: 1,
            Count: 9,
            Plane: "y",
            Center: [3]int{0, 0, 0},
            OutputCount: 2,
            InputFunc: func(nodes []*brain.Node, influences map[string]*brain.Output) {
                for _, node := range nodes {
                    node.Value = 1
                }
            },
        },
    }, false, false)

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
    
    myNet.GenerateAnim(frames, directory)

    fmt.Print("\nSave state?  Enter a name if you wish to save the state:  ")
    fileName, _ := reader.ReadString('\n')
    fileName = strings.TrimSpace(fileName)
    if fileName != "" {
        myNet.SaveState(fileName, directory)
    }
}