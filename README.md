![CircleCI status](https://circleci.com/gh/Firedrake969/gopher-brain.svg?style=shield)

[![Join the chat at https://gitter.im/Firedrake969/gopher-brain](https://badges.gitter.im/Firedrake969/gopher-brain.svg)](https://gitter.im/Firedrake969/gopher-brain?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

### Latest animation
(that I've bothered to upload)

![Latest image](/latest.gif)

### How to run

- `cd` into `example/`
- `go build`
- `./example`
- Enter a string if you'd like to load a preexisting state, or leave the prompt blank if you want to create a fresh state.
- Hit space whenever you'd like to end the simulation
- Enter a string if you'd like to save the current state, or leave the prompt blank if you don't want to save the current state.
- `cd ..`
- `python display_net.py X`
- `cd` into the `frames` directory
- `ffmpeg -framerate Y -i net_%01d.png anim.gif`, Y being fps
- to convert to a video, `ffmpeg -f gif -i anim.gif anim.mp4`
- to make a video directly, `ffmpeg -framerate Y -i net_%01d.png anim.mp4`

Modify `example/.go` to create sensors for the state of the net/brain/whatever.

### Cross-compiling

http://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5

`env GOOS=<your os> GOARCH=<your architecture> go build`

https://golang.org/doc/install/source#environment