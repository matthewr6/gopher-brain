package main

import (
    "strconv"
)

func StrsToInts(strings []string) []int {
    retval := []int{}
    for _, i := range strings {
        j, _ := strconv.Atoi(i)
        retval = append(retval, j)
    }
    return retval
}