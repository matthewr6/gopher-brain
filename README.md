# stuff part 1

note - "phases" for growth
see page 2 of https://pdfs.semanticscholar.org/a4af/fbdf337df0aef29f68364bf9f46338576eaa.pdf

no phase for neuron growth, probably, but phases for higher synaptic creation/deletions

neuron death just means it doesn't output to anything - `suggesting
that a result of the link pruning is a reduction in the
number of effective neurons in the system.`

```
1.) Rapid growth of the number of neurons before birth.
2.) Rapid growth of the number of axons before birth.
3.) Slow synaptogenesis which starts before birth
and continues until puberty.
4.) Very rapid axon loss starting at birth.
5.) Slow apoptose of neurons until puberty.
6.) Slow synapse elimination starting in the
middle of the critical phases and continuing
through the remaining life

Two essential principles of this algorithm are *timing* and
*initial overproduction* and *subsequent elimination*.
```
that's three but whatever.

synapse elimination --> decline in memory/skills?

# stuff that's more important
![CircleCI status](https://circleci.com/gh/matthewr6/gopher-brain.svg?style=shield)

[![Join the chat at https://gitter.im/matthewr6/gopher-brain](https://badges.gitter.im/matthewr6/gopher-brain.svg)](https://gitter.im/matthewr6/gopher-brain?utm_source=badge&utm_medium=badge&utm_campaign=pr-badge&utm_content=badge)

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