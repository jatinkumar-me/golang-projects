package main

import (
	"fmt"
	"strconv"
	"strings"
)

type IntSlice []int

func (intSlice IntSlice) String() string {
    var strs []string

    for _, v := range intSlice {
        strs = append(strs, strconv.Itoa(v))
    }

    return "[" + strings.Join(strs, ";") + "]"
}

func main() {
    var v IntSlice = []int{1, 2, 343, 53}
    fmt.Printf("%T %[1]v\n", v)
}
