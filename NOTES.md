potentially - save sensor/output centers (and possibly nodes?) but not the functions, and have user input functions?  that would be after this branch is merged though.

go through all of my `// todo - ` comments and clean up anything unnecessary

https://www.raspberrypi.org/forums/viewtopic.php?f=29&t=109587

### working on

exploring the possibility of saving centers(and nodes?  speed diff?) of outputs/sensors so you can reset the functions

right now, output functions are just additive, so that wouldn't be a big deal

the output connections *are* unique to the output 

would have to check the same node coordinates in save_test.go; however, would not check whether function esists

also display existing sensors (not outputs) to user on load so that user can cancel/restart based on that information

after custom functions are loaded, runs cript to prune sensors (and their corresponding outputs) that don't have their functions set

as of commit `2411ad5b16`:
- allow users to set custom functions on loaded sensors
- prune unused sensors and their corresponding outputs (after any user custom stuff is set)
- make tests more in depth

for the first two - use regex
sensor format is always `<name>-(one|two)`
output format is always `<name>-(one|two)-<d>` where `d` is the number of the sensor

`[a-zA-Z]+-(one|two)`
`[a-zA-Z]+-(one|two)-(\d+)`

would regex be needed, or can I just split the string on hyphens and grab the first item of the list?

for both, look specifically for the first part - the `name`