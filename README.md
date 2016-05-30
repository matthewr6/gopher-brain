### latest animation

[![Join the chat at https://gitter.im/Firedrake969/gopher-brain](https://badges.gitter.im/Firedrake969/gopher-brain.svg)](https://gitter.im/Firedrake969/gopher-brain?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

![Latest image](/latest.gif)

### how to make things happen

- `cd` into the root of the project
- `go build`
- `gopher-brain`
- Enter a string if you'd like to load a preexisting state, or leave the prompt blank if you want to create a fresh state.
- Hit space whenever you'd like to end the simulation
- Enter a string if you'd like to save the current state, or leave the prompt blank if you don't want to save the current state.
- `python display_net.py X`
- `cd` into the `frames` directory
- `ffmpeg -framerate Y -i net_%01d.png anim.gif`, Y being fps

Modify `main.go` to create sensors for the state of the net/brain/whatever

#### personal notes/links
http://biology.stackexchange.com/questions/5306/how-do-neurons-form-new-connections-in-brain-plasticity