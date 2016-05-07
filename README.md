### to produce an animation

- `go build`
- `gopher-brain X`, X being however many cycles/frames
- `python display_net.py X`
- `ffmpeg -framerate Y -i net_%01d.png anim.gif`, Y being fps

Modify `main.go` to save, load, and stimulate the state of the net/brain/whatever

`SaveState` - saves state to a JSON file
`LoadState` - loads state from a JSON file