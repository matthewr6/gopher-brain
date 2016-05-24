package main

import (
    "strconv"
)

func StrsToInts(strings []string) []int {
    retval := []int{}
    for _, i := range strings {
        j, err := strconv.Atoi(i)
        if err != nil {
            panic(err)
        }
        retval = append(retval, j)
    }
    return retval
}